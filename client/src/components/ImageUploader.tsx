import Image from "next/image"
import React, { ButtonHTMLAttributes, ChangeEvent, memo, useRef } from "react"

interface ImageUploader extends ButtonHTMLAttributes<HTMLButtonElement> {
    src?: string
    fallBackSrc: string
    alt?: string
    width: number
    height: number
    onChangeHandler: (e: ChangeEvent<HTMLInputElement>) => void
}

const ImageUploader = (props: ImageUploader) => {
    const {
        src,
        fallBackSrc,
        alt,
        width,
        height,
        onChangeHandler,
        className,
        ...restProps
    } = props

    // Hooks
    const inputRef = useRef<HTMLInputElement | null>(null)

    // Button Handler
    const btnHandler = () => {
        inputRef.current?.click()
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
                onChange={onChangeHandler}
                className="hidden"
            />
            <Image
                src={src ?? fallBackSrc}
                alt={alt ?? "upload image"}
                width={width}
                height={height}
                className={className}
            />
        </button>
    )
}

export default memo(ImageUploader)
