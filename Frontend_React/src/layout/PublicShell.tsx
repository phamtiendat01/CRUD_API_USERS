export default function PublicShell({children}:{children:React.ReactNode}) {
  return (
    <div className="min-h-screen grid place-items-center p-4">
      <div className="w-full max-w-md card p-6">{children}</div>
    </div>
  );
}
