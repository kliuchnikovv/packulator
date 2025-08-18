import React, { useState } from "react";
import "./App.css";
import PackSetsList from "./components/PackSetsList";
import AddPackSet from "./components/AddPackSet";

interface CalculationResult {
  [size: string]: number;
}

interface PackItem {
  id: string;
  pack_id: string;
  size: number;
}

interface PackSet {
  id: string;
  version_hash: string;
  total_amount: number;
  pack_items: PackItem[];
  created_at: string;
  updated_at: string;
}

function App() {
  const [orderSize, setOrderSize] = useState<string>("");
  const [result, setResult] = useState<CalculationResult | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [selectedPackSet, setSelectedPackSet] = useState<PackSet | null>(null);
  const [activeTab, setActiveTab] = useState<"calculator" | "manage">("calculator");
  const [refreshTrigger, setRefreshTrigger] = useState<number>(0);

  const calculatePacks = async () => {
    if (!orderSize || parseInt(orderSize) <= 0) {
      setError("Please enter a valid order size");
      return;
    }

    if (!selectedPackSet) {
      setError("Please select a pack set");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `http://localhost:8080/packaging/number_of_packages?amount=${orderSize}&packs_hash=${selectedPackSet.version_hash}`,
        {
          method: "GET",
        },
      );

      if (!response.ok) {
        throw new Error("Failed to calculate packs");
      }

      const data = await response.json();
      setResult(data);
      console.log(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const handlePackSetSelect = (packSet: PackSet) => {
    setSelectedPackSet(packSet);
    setActiveTab("calculator");
  };

  const handlePackSetAdded = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="App">
      <div className="app-header">
        <h1>Packulator</h1>
        <div className="tabs">
          <button 
            className={activeTab === "calculator" ? "tab active" : "tab"}
            onClick={() => setActiveTab("calculator")}
          >
            Calculator
          </button>
          <button 
            className={activeTab === "manage" ? "tab active" : "tab"}
            onClick={() => setActiveTab("manage")}
          >
            Manage Pack Sets
          </button>
        </div>
      </div>

      {activeTab === "calculator" && (
        <div className="calculator-container">
          <h2>Order Pack Calculator</h2>

          {selectedPackSet && (
            <div className="selected-pack-set">
              <strong>Using pack set:</strong> #{selectedPackSet.id.substring(0, 8)}
              <span className="pack-sizes">({selectedPackSet.pack_items.map(item => item.size).sort((a, b) => a - b).join(", ")})</span>
            </div>
          )}

          <div className="input-section">
            <label htmlFor="orderSize">Order Size:</label>
            <input
              id="orderSize"
              type="number"
              value={orderSize}
              onChange={(e) => setOrderSize(e.target.value)}
              placeholder="Enter number of items"
              min="1"
            />
            <button onClick={calculatePacks} disabled={loading || !selectedPackSet}>
              {loading ? "Calculating..." : "Calculate Packs"}
            </button>
          </div>

          {!selectedPackSet && (
            <div className="warning-message">
              Please go to "Manage Pack Sets" tab to select a pack set first.
            </div>
          )}

          {error && <div className="error-message">{error}</div>}

          {result && (
            <div className="result-section">
              <h3>Calculation Result</h3>
              <div className="total-packs">
                <strong>Total Packs Needed: {Object.values(result).reduce((sum, count) => sum + count, 0)}</strong>
              </div>

              <div className="pack-breakdown">
                <h4>Pack Breakdown:</h4>
                <ul>
                  {Object.entries(result)
                    .sort(([a], [b]) => parseInt(b) - parseInt(a))
                    .map(([size, count]) => (
                      <li key={size}>
                        Pack size {size}: {count} packs
                      </li>
                    ))}
                </ul>
              </div>
            </div>
          )}
        </div>
      )}

      {activeTab === "manage" && (
        <div className="manage-container">
          <div className="manage-sections">
            <div className="add-section">
              <AddPackSet onAdd={handlePackSetAdded} />
            </div>
            <div className="list-section">
              <PackSetsList 
                onSelect={handlePackSetSelect}
                refreshTrigger={refreshTrigger}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
