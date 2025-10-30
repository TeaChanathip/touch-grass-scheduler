import Image from "next/image"
import React, { ButtonHTMLAttributes, useRef } from "react"

interface ImageUploader extends ButtonHTMLAttributes<HTMLButtonElement> {
    src?: string
    fallBackSrc: string
    alt?: string
    width: number
    height: number
}

export default function ImageUploader(props: ImageUploader) {
    const { src, fallBackSrc, alt, width, height, className, ...restProps } =
        props

    // Hooks
    const inputRef = useRef<HTMLInputElement | null>(null)

    // Button Handler
    const btnHandler = () => {
        inputRef.current?.click()
    }

    const inputChangeHandler = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0]
        console.log(file)
    }

    return (
        <button
            type="button"
            onClick={btnHandler}
            {...restProps}
            className={`cursor-pointer disabled:pointer-events-none ${className}`}
        >
            <input
                type="file"
                accept="image/png, image/jpeg, image/webp"
                ref={inputRef}
                onChange={inputChangeHandler}
                className="hidden"
            />
            <Image
                src={src ?? fallBackSrc}
                alt={alt ?? "upload image"}
                width={width}
                height={height}
            />
        </button>
    )
}
