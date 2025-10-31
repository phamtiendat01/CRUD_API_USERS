export default function Chip({
  children, tone="gray", className=""
}: { children: React.ReactNode; tone?: "brand"|"gray"|"green"|"red"|"amber"; className?: string; }) {
  const map = {
    gray:  "bg-gray-100 text-gray-700",
    brand: "bg-brand-100 text-brand-700",
    green: "bg-green-100 text-green-700",
    red:   "bg-red-100 text-red-700",
    amber: "bg-amber-100 text-amber-700",
  } as const;
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${map[tone]} ${className}`}>{children}</span>
  );
}
