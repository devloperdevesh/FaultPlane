import asyncio
import os
import platform
import socket
import struct
import time
from typing import Dict, Optional, Tuple, NamedTuple, Callable

SOL_TCP = getattr(socket, "SOL_TCP", 6)
TCP_INFO = getattr(socket, "TCP_INFO", 11)

# Struct layout for Linux tcp_info (extracting essential RTT & window fields)
# Fields extracted: state, rto, snd_mss, rcv_mss, unacked, lost, retrans, rtt, rttvar, snd_cwnd, rcv_space, total_retrans
TCP_INFO_FMT = "BBBBBBBBIIIIIIIIIIIIIIII"
TCP_INFO_SIZE = struct.calcsize(TCP_INFO_FMT)


class TCPKernelMetrics(NamedTuple):
    """Extracted Linux kernel TCP metric state."""
    rtt_us: int           # Smoothed Round-Trip Time in microseconds
    rttvar_us: int        # RTT Variance (Jitter) in microseconds
    snd_cwnd: int         # Congestion Window size in packets
    rcv_space: int        # Current receive buffer space estimate in bytes
    unacked: int          # Number of unacknowledged packets
    lost: int             # Number of lost packets
    retransmits: int      # Retransmission counter
    total_retrans: int    # Cumulative retransmissions


class WindowOptimizationResult(NamedTuple):
    """Computed target buffer window sizing and congestion risk score."""
    target_rcvbuf: int
    target_sndbuf: int
    jitter_factor: float
    risk_score: float
    clamped: bool


class TCPKernelProfiler:
    """Reads socket buffer metadata directly from the Linux kernel TCP stack."""

    @staticmethod
    def get_tcp_info(sock: socket.socket) -> Optional[TCPKernelMetrics]:
        """Reads struct tcp_info from the provided TCP socket."""
        if platform.system() != "Linux":
            # Fallback/stub for non-Linux environments during dev/testing
            return TCPKernelMetrics(
                rtt_us=500, rttvar_us=50, snd_cwnd=10, rcv_space=65536,
                unacked=0, lost=0, retransmits=0, total_retrans=0
            )

        try:
            raw_info = sock.getsockopt(SOL_TCP, TCP_INFO, 104)
            if len(raw_info) < 92:  # Ensure buffer contains minimal header fields
                return None
            
            # Unpack key header fields from tcp_info
            data = struct.unpack("8B10I", raw_info[:48])
            rtt = data[16]          # tcpi_rtt
            rttvar = data[17]       # tcpi_rttvar
            snd_cwnd = data[19]     # tcpi_snd_cwnd
            rcv_space = data[20]    # tcpi_rcv_space
            unacked = data[12]      # tcpi_unacked
            lost = data[14]         # tcpi_lost
            retransmits = data[15]  # tcpi_retrans
            total_retrans = data[21] if len(raw_info) >= 92 else 0

            return TCPKernelMetrics(
                rtt_us=rtt,
                rttvar_us=rttvar,
                snd_cwnd=snd_cwnd,
                rcv_space=rcv_space,
                unacked=unacked,
                lost=lost,
                retransmits=retransmits,
                total_retrans=total_retrans
            )
        except (OSError, struct.error):
            return None


class TCPWindowOptimizationMatrix:
    """Calculates optimal socket buffer sizes dynamically based on RTT variance and BDP."""

    def __init__(self, 
                 min_buffer_size: int = 64 * 1024,      # 64 KB
                 max_buffer_size: int = 16 * 1024 * 1024, # 16 MB
                 alpha_ema: float = 0.2,
                 jitter_threshold_us: float = 2000.0):   # 2 ms jitter threshold
        self.min_buffer_size = min_buffer_size
        self.max_buffer_size = max_buffer_size
        self.alpha_ema = alpha_ema
        self.jitter_threshold_us = jitter_threshold_us
        
        self._ema_rttvar: Dict[int, float] = {}

    def compute_optimal_windows(self, sock_fd: int, metrics: TCPKernelMetrics) -> WindowOptimizationResult:
        """Computes new TCP send and receive buffer parameters for a socket."""
        # Update Exponential Moving Average (EMA) of RTT Variance (Jitter)
        prev_ema = self._ema_rttvar.get(sock_fd, float(metrics.rttvar_us))
        curr_ema = (self.alpha_ema * metrics.rttvar_us) + ((1.0 - self.alpha_ema) * prev_ema)
        self._ema_rttvar[sock_fd] = curr_ema

        # Calculate Jitter Factor relative to baseline target
        jitter_factor = curr_ema / max(self.jitter_threshold_us, 1.0)
        
        # Calculate Congestion Risk Score based on unacked packets, lost packets, and jitter
        risk_score = min(1.0, (metrics.lost * 0.4) + (metrics.unacked * 0.1) + (jitter_factor * 0.5))

        rtt_sec = max(metrics.rtt_us, 1) / 1_000_000.0
        # Estimated throughput capability based on current MSS & window space
        est_bw_bytes_sec = max(metrics.rcv_space, metrics.snd_cwnd * 1460) / max(rtt_sec, 0.0001)
        base_bdp = int(est_bw_bytes_sec * rtt_sec)

        if risk_score > 0.7:
            # High risk: Clamp buffer growth to prevent queue bloat and drop bursts
            scale_factor = 1.0
            clamped = True
        elif jitter_factor > 1.5:
            # High jitter: Moderately increase buffer to absorb latency spikes
            scale_factor = 1.5 + min(jitter_factor, 2.0)
            clamped = False
        else:
            # Optimal low-latency state: Keep buffers tight for microsecond latency
            scale_factor = 1.2
            clamped = False

        target_rcv = int(max(self.min_buffer_size, min(base_bdp * scale_factor, self.max_buffer_size)))
        target_snd = int(max(self.min_buffer_size, min(target_rcv * 1.2, self.max_buffer_size)))

        return WindowOptimizationResult(
            target_rcvbuf=target_rcv,
            target_sndbuf=target_snd,
            jitter_factor=jitter_factor,
            risk_score=risk_score,
            clamped=clamped
        )

    def cleanup_socket(self, sock_fd: int) -> None:
        """Removes state tracking for closed sockets."""
        self._ema_rttvar.pop(sock_fd, None)


class AsyncTCPWindowOptimizer:
    """Asynchronous background manager for dynamic kernel socket buffer optimization."""

    def __init__(self, 
                 poll_interval_sec: float = 0.05, 
                 telemetry_cb: Optional[Callable[[int, TCPKernelMetrics, WindowOptimizationResult], None]] = None):
        self.poll_interval_sec = poll_interval_sec
        self.telemetry_cb = telemetry_cb
        self.matrix = TCPWindowOptimizationMatrix()
        self._monitored_sockets: Dict[int, socket.socket] = {}
        self._task: Optional[asyncio.Task] = None
        self._running = False

    def register_socket(self, sock: socket.socket) -> None:
        """Registers a TCP socket for asynchronous optimization profiling."""
        fd = sock.fileno()
        if fd != -1:
            self._monitored_sockets[fd] = sock

    def unregister_socket(self, sock: socket.socket) -> None:
        """Unregisters a TCP socket from optimization."""
        fd = sock.fileno()
        if fd in self._monitored_sockets:
            del self._monitored_sockets[fd]
            self.matrix.cleanup_socket(fd)

    async def start() -> None:
        """Starts the optimization background daemon loop."""
        self._running = True
        self._task = asyncio.create_task(self._run_optimization_loop())

    async def stop() -> None:
        """Stops the optimization background daemon loop."""
        self._running = False
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except asyncio.CancelledError:
                pass

    async def _run_optimization_loop() -> None:
        """Non-blocking asynchronous loop executing passive RTT kernel profiling."""
        while self._running:
            for fd, sock in list(self._monitored_sockets.items()):
                try:
                    metrics = TCPKernelProfiler.get_tcp_info(sock)
                    if metrics is None:
                        continue

                    res = self.matrix.compute_optimal_windows(fd, metrics)

                    # Dynamic socket buffer adjustments via setsockopt
                    sock.setsockopt(socket.SOL_SOCKET, socket.SO_RCVBUF, res.target_rcvbuf)
                    sock.setsockopt(socket.SOL_SOCKET, socket.SO_SNDBUF, res.target_sndbuf)

                    if self.telemetry_cb:
                        self.telemetry_cb(fd, metrics, res)

                except (OSError, ValueError):
                    # Handle socket closure or invalid FDs
                    self.unregister_socket(sock)

            await asyncio.sleep(self.poll_interval_sec)


# --- Example Usage & Self-Test ---
async def main():
    print("Initializing Asynchronous TCP Window Optimization Engine...")
    
    def on_telemetry(fd: int, metrics: TCPKernelMetrics, result: WindowOptimizationResult):
        print(f"[FD {fd}] RTT: {metrics.rtt_us}us | Jitter: {metrics.rttvar_us}us | "
              f"Risk Score: {result.risk_score:.2f} | RCVBUF: {result.target_rcvbuf} B | Clamped: {result.clamped}")

    optimizer = AsyncTCPWindowOptimizer(poll_interval_sec=0.1, telemetry_cb=on_telemetry)
    
    # Create a local test socket pair
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind(("127.0.0.1", 0))
    server.listen(1)
    
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(server.getsockname())
    
    conn, _ = server.accept()

    # Register active client & connection sockets
    optimizer.register_socket(client)
    optimizer.register_socket(conn)

    optimizer_task = asyncio.create_task(optimizer._run_optimization_loop())
    await asyncio.sleep(0.3)
    
    optimizer._running = False
    optimizer_task.cancel()
    client.close()
    conn.close()
    server.close()
    print("Optimization engine run completed successfully.")

if __name__ == "__main__":
    asyncio.run(main())