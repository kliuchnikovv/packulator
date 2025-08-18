import React, { useState } from "react";

interface AddPackSetProps {
  onAdd?: () => void;
}

const AddPackSet: React.FC<AddPackSetProps> = ({ onAdd }) => {
  const [packsInput, setPacksInput] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!packsInput.trim()) {
      setError("Please enter pack sizes");
      return;
    }

    const packsArray = packsInput
      .split(",")
      .map((p) => p.trim())
      .filter((p) => p !== "")
      .map((p) => parseInt(p))
      .filter((p) => !isNaN(p) && p > 0);

    if (packsArray.length === 0) {
      setError(
        "Please enter valid pack sizes (positive numbers separated by commas)",
      );
      return;
    }

    const uniquePacks = Array.from(new Set(packsArray)).sort((a, b) => a - b);

    setLoading(true);
    setError("");
    setSuccess("");

    try {
      const response = await fetch("http://localhost:8080/packs/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          packs: uniquePacks,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Failed to create pack set");
      }

      setSuccess(
        `Pack set with sizes [${uniquePacks.join(", ")}] created successfully!`,
      );
      setPacksInput("");

      if (onAdd) {
        onAdd();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  const clearMessages = () => {
    setError("");
    setSuccess("");
  };

  return (
    <div className="add-pack-set">
      <h2>Add New Pack Set</h2>

      <form onSubmit={handleSubmit} className="pack-set-form">
        <div className="form-group">
          <label htmlFor="packs">Pack Sizes:</label>
          <input
            id="packs"
            type="text"
            value={packsInput}
            onChange={(e) => {
              setPacksInput(e.target.value);
              clearMessages();
            }}
            placeholder="e.g., 250, 500, 1000, 2000, 5000"
            required
          />
          <small className="help-text">
            Enter pack sizes separated by commas. Duplicates will be removed.
          </small>
        </div>

        <button type="submit" disabled={loading} className="submit-btn">
          {loading ? "Creating..." : "Create Pack Set"}
        </button>
      </form>

      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}

      <div className="preview">
        <h3>Preview:</h3>
        <p>
          <strong>Pack sizes:</strong>{" "}
          {packsInput
            ? Array.from(
                new Set(
                  packsInput
                    .split(",")
                    .map((p) => p.trim())
                    .filter((p) => p !== "")
                    .map((p) => parseInt(p))
                    .filter((p) => !isNaN(p) && p > 0),
                ),
              )
                .sort((a, b) => a - b)
                .join(", ")
            : "Not set"}
        </p>
        <p>
          <strong>Total Amount:</strong>{" "}
          {packsInput
            ? Math.max(
                0,
                ...packsInput
                  .split(",")
                  .map((p) => p.trim())
                  .filter((p) => p !== "")
                  .map((p) => parseInt(p))
                  .filter((p) => !isNaN(p) && p > 0),
              ) || 0
            : 0}
        </p>
      </div>
    </div>
  );
};

export default AddPackSet;
