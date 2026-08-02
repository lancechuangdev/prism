/* eslint-disable react-refresh/only-export-components -- routing primitives intentionally share one dependency-free module */
import {
  useEffect,
  useState,
  type AnchorHTMLAttributes,
  type MouseEvent,
} from 'react'

const navigationEvent = 'prism:navigate'

export function navigate(to: string) {
  window.history.pushState({}, '', to)
  window.dispatchEvent(new Event(navigationEvent))
}

export function usePathname() {
  const [pathname, setPathname] = useState(window.location.pathname)

  useEffect(() => {
    const updatePathname = () => setPathname(window.location.pathname)
    window.addEventListener('popstate', updatePathname)
    window.addEventListener(navigationEvent, updatePathname)
    return () => {
      window.removeEventListener('popstate', updatePathname)
      window.removeEventListener(navigationEvent, updatePathname)
    }
  }, [])

  return pathname
}

type LinkProps = Omit<AnchorHTMLAttributes<HTMLAnchorElement>, 'href'> & {
  to: string
}

export function Link({ to, onClick, ...props }: LinkProps) {
  function handleClick(event: MouseEvent<HTMLAnchorElement>) {
    onClick?.(event)
    if (
      event.defaultPrevented ||
      event.button !== 0 ||
      event.metaKey ||
      event.ctrlKey ||
      event.shiftKey ||
      event.altKey
    )
      return
    event.preventDefault()
    navigate(to)
  }

  return <a href={to} onClick={handleClick} {...props} />
}
