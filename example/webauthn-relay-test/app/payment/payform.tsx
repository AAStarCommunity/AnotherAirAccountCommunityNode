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
    const [networks, setNetworks] = useState<{id: string; name: string}[]>([]);
    const [selectedNetwork, setSelectedNetwork] = useState<string>('');
    
    useEffect(() => {
      api.get(
        API.SUPPORT_NETWORKS,
      ).then(response => {
        const data: { [key: string]: boolean } = response.data.data;
        const networksArray = Object.keys(data)
        .filter(key => data[key] === true).map(key => ({
          id: key,
          name: key,
        }));
        setNetworks(networksArray);
      })
    }, []);

    const handleNetworkChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
      setSelectedNetwork(event.target.value);
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
        </div>
        {children}
      </form>
    );
  }
  