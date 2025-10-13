import { ButtonHTMLAttributes } from "react"
import "../styles/MyButton.css"

interface MyButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant: "positive" | "neutral" | "negative"
}

export default function MyButton(props: MyButtonProps) {
    const { variant, children, className, ...restProps } = props

    return (
        <button
            {...restProps}
            className={`text-2xl py-2 px-5 rounded-xl cursor-pointer 
                ${variant} ${className ?? ""}`}
        >
            {children}
        </button>
    )
}
