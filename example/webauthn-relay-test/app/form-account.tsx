export function FormAccount({
  action,
  children,
  isDiscoverable,
}: {
  action: any;
  children: React.ReactNode;
  isDiscoverable: boolean;
}) {
  return (
    <form
      action={action}
      className="flex flex-col space-y-4 bg-gray-50 px-4 py-8 sm:px-16"
    >
      {!isDiscoverable && (
        <div>
          <label
            htmlFor="zuzalu-city"
            className="block text-xs text-gray-600 uppercase"
          >
            Zuzalu City Unique ID
          </label>
          <input
            id="zuzalu-city"
            name="zuzalu-city"
            placeholder="Your Zuzalu City Unique ID"
            autoComplete="zuzalu-city"
            required
            className="mt-1 block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-black focus:outline-none focus:ring-black sm:text-sm"
          />
        </div>
      )}
      {children}
    </form>
  );
}
