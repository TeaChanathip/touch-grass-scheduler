"use client"

import { VisibilityOffRounded, VisibilityRounded } from "@mui/icons-material"
import { CSSProperties, InputHTMLAttributes, memo, useState } from "react"
import { get, useFormContext } from "react-hook-form"
import StatusMessage from "./StatusMessage"

interface FormPasswordProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string
    name: string
    warn?: boolean
    hideMsg?: boolean
}

const FormPassword = ({
    label,
    name,
    warn,
    hideMsg,
    ...restProps
}: FormPasswordProps) => {
    // Hooks
    const [isShowPassword, setShowPassword] = useState(false)

    // Form Context
    const {
        register,
        formState: { errors },
    } = useFormContext()
    const warningMsg = errors[name]?.message as string | undefined

    const toggleShowPassword = () => {
        setShowPassword(!isShowPassword)
    }

    const constructAdditionalStyle = () => {
        const style: CSSProperties = {}
        if (warn || warningMsg) {
            style.borderColor = "var(--color-prim-red)"
        }
        style.paddingRight = "40px"
        return style
    }

    return (
        <div className="flex flex-col">
            {label && (
                <label className="text-2xl">
                    {label}
                    {restProps.required && (
                        <span className="text-prim-red">*</span>
                    )}
                </label>
            )}
            <span className="flex flex-row items-center justify-end">
                <input
                    type={isShowPassword ? "text" : "password"}
                    {...restProps}
                    {...register(name)}
                    className="w-full h-11 px-3 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                    style={constructAdditionalStyle()}
                />
                <span
                    className="absolute mr-3 cursor-pointer"
                    onClick={() => toggleShowPassword()}
                >
                    {isShowPassword ? (
                        <VisibilityOffRounded />
                    ) : (
                        <VisibilityRounded />
                    )}
                </span>
            </span>
            {!hideMsg && (
                <StatusMessage
                    msg={warningMsg}
                    variant="error"
                    className="text-xl"
                />
            )}
        </div>
    )
}

export default memo(FormPassword)
