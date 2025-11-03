import { RefObject, useEffect } from "react"

export default function useClickOutside(
    ref: RefObject<any>,
    onClickOutside: () => void
) {
    useEffect(() => {
        function handleClickOutside(event: Event) {
            if (ref.current && !ref.current.contains(event.target)) {
                onClickOutside()
            }
        }

        document.addEventListener("mousedown", handleClickOutside)

        return () =>
            document.removeEventListener("mousedown", handleClickOutside)
    }, [ref, onClickOutside])
}
