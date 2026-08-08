import RoutePage from "@/components/layout/RoutePage";
import WasmSandbox from "@/components/experimental/wasm/WasmSandbox";

export default function WasmPage() {
  return (
    <RoutePage
      eyebrow="Experiment"
      title="WASM Sandbox"
      description="Run isolated runtime experiments in the WebAssembly sandbox."
    >
      <WasmSandbox />
    </RoutePage>
  );
}
