import { useEffect, useState } from "react";
import API from "../api/api";
import api from "../api";

export function PayForm({
  action,
  children,
}: {
  action: any;
  children: React.ReactNode;
}) {
  const [networks, setNetworks] = useState<{ id: string; name: string }[]>([]);
  const [selectedNetwork, setSelectedNetwork] = useState<string>('');

  useEffect(() => {
    api.get(API.SUPPORT_NETWORKS).then(response => {
      const data: { [key: string]: boolean } = response.data.data;
      const networksArray = Object.keys(data)
        .filter(key => data[key] === true).map(key => ({
          id: key,
          name: key,
        }));
      setNetworks(networksArray);
      setSelectedNetwork(networksArray[0].id);
    })
  }, []);

  const handleNetworkChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedNetwork(event.target.value);
  };

  const handleAddNetwork = async () => {
    if (!selectedNetwork) return;
    try {
      const response = await api.post(
        API.CREATE_AA,
        {
          network: selectedNetwork
        },
        {
          headers: {
            Authorization: "Bearer " + localStorage.getItem("token"),
          },
        }
      );
      alert("create aa success");
    } catch (error) {
      console.error('Error adding network:', error);
    }
  };

  return (
    <form
      action={action}
      className="flex flex-col space-y-4 bg-gray-50 px-4 py-8 sm:px-16"
    >
      <div>
        <label
          htmlFor="txdata"
          className="block text-xs text-gray-600 uppercase"
        >
          TxData
        </label>
        <input
          id="txdata"
          name="txdata"
          type="txdata"
          value="48656c6c6f2c20576f726c6421"
          placeholder="Hashed UserOp"
          required
          className="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-black focus:outline-none focus:ring-black sm:text-sm"
        />
        <label
          htmlFor="network"
          className="block text-xs text-gray-600 uppercase mt-4"
        >
          Network
        </label>
        <div className="flex gap-2">
          <select
            value={selectedNetwork}
            onChange={handleNetworkChange}
            title="network"
            id="network"
            name="network"
            className="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-black focus:outline-none focus:ring-black sm:text-sm"
          >
            {networks.map(network => (
              <option key={network.id} value={network.id}>
                {network.name}
              </option>
            ))}
          </select>
          <button
            type="button"
            onClick={handleAddNetwork}
            className="mt-1 inline-flex items-center justify-center rounded-md border border-gray-300 px-3 py-2 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
      {children}
    </form>
  );
}