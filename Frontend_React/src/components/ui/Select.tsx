export default function Select(
  { className = "", ...rest }: React.SelectHTMLAttributes<HTMLSelectElement>
) {
  return (
    <select
      className={`input pr-8 ${className}`}
      {...rest}
    />
  );
}
