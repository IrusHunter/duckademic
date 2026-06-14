// Іконки через той самий SVG-спрайт, що й у web: /img/icons.svg#<id>.
// Спрайт лежить у public/img/icons.svg (shell і home-app).
interface IconProps {
  id: string
  size?: number
  width?: number
  height?: number
  className?: string
}

export default function Icon({ id, size = 16, width, height, className }: IconProps) {
  return (
    <svg
      width={width ?? size}
      height={height ?? size}
      className={className}
      aria-hidden="true"
    >
      <use href={`/img/icons.svg#${id}`} />
    </svg>
  )
}
