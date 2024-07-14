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
            htmlFor="amount"
            className="block text-xs text-gray-600 uppercase"
          >
            Transfer Amount
          </label>
          <input
            id="amount"
            name="amount"
            type="amount"
            placeholder="$1.00"
            required
            className="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-black focus:outline-none focus:ring-black sm:text-sm"
          />
        </div>
        {children}
      </form>
    );
  }
  