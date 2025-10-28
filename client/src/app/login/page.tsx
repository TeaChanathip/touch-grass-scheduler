"use client"

import { useForm, UseFormHandleSubmit } from "react-hook-form"
import FormStringInput from "../../components/FormStringInput"
import MyButton from "../../components/MyButton"
import Link from "next/link"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { useAppDispatch, useAppSelector } from "../../store/hooks"
import {
    userLogin,
    selectUserErrMsg,
    selectUserStatus,
} from "../../store/features/user/userSlice"
import { useEffect } from "react"
import { useRouter } from "next/navigation"
import StatusMessage from "../../components/StatusMessage"

export default function LoginPage() {
    // Store
    const userStatus = useAppSelector(selectUserStatus)

    // Hooks
    const router = useRouter()

    useEffect(() => {
        if (userStatus === "authenticated") {
            router.push("/")
        }
    }, [userStatus, router])

    return <LoginForm />
}

// Form Schema
const schema = z.object({
    email: z.email("Invalid email format"),
    password: z
        .string()
        .min(8, "At least 8 characters")
        .max(64, "At most 64 characters"),
})

function LoginForm() {
    // Hooks
    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    return (
        <form className="w-[70vw] max-w-96 flex flex-col gap-5">
            <FormStringInput
                label="Email Address"
                type="email"
                required
                register={register("email")}
                warn={valErrors.email !== undefined}
                warningMsg={valErrors.email?.message}
            />
            <FormStringInput
                label="Password"
                type="password"
                required
                register={register("password")}
                warn={valErrors.password !== undefined}
                warningMsg={valErrors.password?.message}
            />
            <ButtonSection
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
                hasValidationErr={Object.keys(valErrors).length != 0}
            />
        </form>
    )
}

function ButtonSection({
    handleSubmit,
    isSubmitting,
    hasValidationErr,
}: {
    handleSubmit: UseFormHandleSubmit<z.infer<typeof schema>>
    isSubmitting: boolean
    hasValidationErr: boolean
}) {
    // Store
    const dispatch = useAppDispatch()
    const userErrMsg = useAppSelector(selectUserErrMsg)

    // Hooks
    const router = useRouter()

    const submitHandler = async (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)
        await dispatch(userLogin(result.data!)) // Garanteed to be valid
    }

    return (
        <section className="flex flex-col items-center">
            <Link href="/forgot-password" className="w-fit mb-3 underline">
                Forgot the password?
            </Link>
            <StatusMessage
                msg={userErrMsg}
                variant="error"
                className="text-2xl"
            />
            <div className="w-full mt-3 flex flex-col md:flex-row gap-5">
                <MyButton
                    variant="positive"
                    type="submit"
                    disabled={hasValidationErr || isSubmitting}
                    onClick={handleSubmit(submitHandler)}
                    className="w-full md:w-44"
                >
                    Login
                </MyButton>
                <MyButton
                    variant="neutral"
                    type="button"
                    className="w-full md:w-44"
                    onClick={() => router.push("/register/verify-email")}
                >
                    Register
                </MyButton>
            </div>
        </section>
    )
}
