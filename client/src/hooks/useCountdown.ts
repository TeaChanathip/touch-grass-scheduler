import { useRef, useState } from "react"

// duration (seconds)
export default function useCountdown(
    duration: number
): [number, () => void, () => void] {
    const [countdown, setCountdown] = useState(0)
    const intervalIdRef = useRef<ReturnType<typeof setInterval>>(undefined)

    const startCountdown = () => {
        // Clear existing interval
        if (intervalIdRef.current) {
            clearInterval(intervalIdRef.current)
        }

        setCountdown(duration)

        intervalIdRef.current = setInterval(() => {
            setCountdown((currentCountDown) => {
                if (currentCountDown <= 1) {
                    clearInterval(intervalIdRef.current)
                    return 0
                }
                return currentCountDown - 1
            })
        }, 1000)
    }

    const stopCountdown = () => {
        if (intervalIdRef.current) {
            clearInterval(intervalIdRef.current)
        }
    }

    return [countdown, startCountdown, stopCountdown]
}
