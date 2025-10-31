export default function Avatar({ src, name, className="" }:{
  src?: string; name?: string; className?: string;
}) {
  const initials = (name || "?")
    .split(" ").map(s=>s[0]).join("").slice(0,2).toUpperCase();
  return src ? (
    <img src={src} alt={name} className={`h-9 w-9 rounded-full object-cover ${className}`} />
  ) : (
    <div className={`h-9 w-9 rounded-full bg-gray-200 grid place-items-center text-xs font-semibold ${className}`}>
      {initials}
    </div>
  );
}
