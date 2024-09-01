export function PayForm({
    action,
    children,
  }: {
    action: any;
    children: React.ReactNode;
  }) {
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
        </div>
        {children}
      </form>
    );
  }
  