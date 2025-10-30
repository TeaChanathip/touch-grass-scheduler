"use client"

import { FormProvider, useForm, useFormContext } from "react-hook-form"
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
import { memo, useEffect } from "react"
import { useRouter } from "next/navigation"
import StatusMessage from "../../components/StatusMessage"
import FormPassword from "../../components/FormPassword"

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

const LoginForm = memo(function LoginForm() {
    // Hooks
    const formMethods = useForm({
        resolver: zodResolver(schema),
        mode: "onChange",
    })

    return (
        <FormProvider {...formMethods}>
            <form className="w-[70vw] max-w-96 flex flex-col gap-5">
                <FormStringInput
                    label="Email Address"
                    type="email"
                    name="email"
                    required
                />
                <FormPassword
                    label="Password"
                    type="password"
                    name="password"
                    required
                />
                <ButtonSection />
            </form>
        </FormProvider>
    )
})

const ButtonSection = memo(function ButtonSection() {
    // Form Context
    const {
        handleSubmit,
        formState: { isValid, isSubmitting },
    } = useFormContext()

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
                    disabled={!isValid || isSubmitting}
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
})
