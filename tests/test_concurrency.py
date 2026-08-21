import pytest
from core.resilience.load_shedder import DynamicConcurrencyLoadShedder

@pytest.mark.asyncio
async def test_load_shedder_preemption_logic():
    shedder = DynamicConcurrencyLoadShedder(critical_latency_threshold_secs=0.100)
    # Simulate high rolling latency anomaly
    shedder.rolling_latency_average = 0.150 
    
    allowed, reason = await shedder.evaluate_transaction({"priority": "NON_CRITICAL"})
    assert allowed is False
    assert reason == "REJECTED_SYSTEM_SATURATED"
