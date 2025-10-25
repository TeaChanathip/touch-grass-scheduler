"use client"

import { useForm } from "react-hook-form"
import FormStringInput from "../../../components/FormStringInput"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import MyButton from "../../../components/MyButton"
import { ApiService } from "../../../services/api.service"
import { AuthService } from "../../../services/auth/auth.service"
import { useRef, useState } from "react"
import StatusMessage from "../../../components/StatusMessage"

export default function VerifyEmailPage() {
    const [responseMsg, setResponseMsg] = useState<
        { msg: string; variant: "success" | "error" } | undefined
    >(undefined)

    // Hooks for button countdown
    const [countDown, setCountdown] = useState(0)
    const intervalIdRef = useRef<ReturnType<typeof setInterval>>(undefined)

    const schema = z.object({
        email: z.email("Invalid email"),
    })

    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    // Submit handler
    const apiService = new ApiService()
    const authService = new AuthService(apiService)
    const submitHandler = async (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)

        try {
            await authService.getRegistrationMail(result.data!.email)
            setResponseMsg({
                msg: "Verification email was sent",
                variant: "success",
            })

            // Clear existing interval
            if (intervalIdRef.current) {
                clearInterval(intervalIdRef.current)
            }

            setCountdown(30)
            intervalIdRef.current = setInterval(() => {
                setCountdown((currentCountDown) => {
                    if (currentCountDown <= 1) {
                        clearInterval(intervalIdRef.current)
                        return 0
                    }
                    return currentCountDown - 1
                })
            }, 1000)
        } catch (err) {
            if (err instanceof Error)
                setResponseMsg({ msg: err.message, variant: "error" })
        }
    }

    return (
        <div className="flex flex-col gap-10 items-center pt-10">
            <h1 className="text-5xl">Verfiy Email</h1>
            <form
                onSubmit={handleSubmit(submitHandler)}
                className="w-[70vw] max-w-96 flex flex-col gap-5"
            >
                <FormStringInput
                    type="email"
                    label="Email Address"
                    required
                    register={register("email")}
                    warn={valErrors.email !== undefined}
                    warningMsg={valErrors.email?.message}
                />
                <div className="w-full flex flex-col items-center">
                    <StatusMessage
                        msg={responseMsg?.msg}
                        variant={responseMsg?.variant ?? "info"}
                        className="text-2xl"
                    />
                    <MyButton
                        variant="positive"
                        type="submit"
                        disabled={
                            Object.keys(valErrors).length > 0 ||
                            isSubmitting ||
                            countDown !== 0
                        }
                        className="w-full md:w-44"
                    >
                        Send{countDown !== 0 && ` (${countDown})`}
                    </MyButton>
                </div>
            </form>
        </div>
    )
}
