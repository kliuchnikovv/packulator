import React, { useState, useEffect } from "react";

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

interface PackSetsListProps {
  onSelect?: (packSet: PackSet) => void;
  onDelete?: (packSet: PackSet) => void;
  refreshTrigger?: number;
}

const PackSetsList: React.FC<PackSetsListProps> = ({
  onSelect,
  onDelete,
  refreshTrigger,
}) => {
  const [packSets, setPackSets] = useState<PackSet[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const fetchPackSets = async () => {
    setLoading(true);
    setError("");

    try {
      const response = await fetch("http://localhost:8080/packs/list", {
        method: "GET",
      });

      if (!response.ok) {
        throw new Error("Failed to fetch pack sets");
      }

      const data = await response.json();
      setPackSets(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPackSets();
  }, [refreshTrigger]);

  const handleDelete = async (packSet: PackSet) => {
    if (
      !window.confirm(
        `Are you sure you want to delete pack set with ID "${packSet.id}"?`,
      )
    ) {
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/packs/delete?id=${packSet.id}`,
        {
          method: "DELETE",
        },
      );

      if (!response.ok) {
        throw new Error("Failed to delete pack set");
      }

      setPackSets(packSets.filter((p) => p.id !== packSet.id));

      if (onDelete) {
        onDelete(packSet);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "An error occurred while deleting",
      );
    }
  };

  if (loading) {
    return <div className="loading">Loading pack sets...</div>;
  }

  if (error) {
    return (
      <div className="error-message">
        {error}
        <button onClick={fetchPackSets}>Retry</button>
      </div>
    );
  }

  return (
    <div className="pack-sets-list">
      <h2>Pack Sets</h2>
      {packSets.length === 0 ? (
        <p>No pack sets available.</p>
      ) : (
        <div className="pack-sets-grid">
          {packSets.map((packSet) => (
            <div key={packSet.id} className="pack-set-item">
              <div className="pack-set-header">
                <h3>Pack Set #{packSet.id.substring(0, 8)}</h3>
                <div className="pack-set-actions">
                  {onSelect && (
                    <button
                      className="select-btn"
                      onClick={() => onSelect(packSet)}
                    >
                      Select
                    </button>
                  )}
                  <button
                    className="delete-btn"
                    onClick={() => handleDelete(packSet)}
                  >
                    Delete
                  </button>
                </div>
              </div>
              <div className="pack-sizes">
                <strong>Pack sizes:</strong>{" "}
                {packSet.pack_items
                  .map((item) => item.size)
                  .sort((a, b) => a - b)
                  .join(", ")}
              </div>
              <div className="pack-info">
                <div>
                  <strong>Total Amount:</strong> {packSet.total_amount}
                </div>
                <div>
                  <strong>Version Hash:</strong> {packSet.version_hash}
                </div>
                <div>
                  <small>
                    Created: {new Date(packSet.created_at).toLocaleString()}
                  </small>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PackSetsList;
