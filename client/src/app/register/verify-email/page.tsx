"use client"

import { useForm, UseFormHandleSubmit } from "react-hook-form"
import FormStringInput from "../../../components/FormStringInput"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import MyButton from "../../../components/MyButton"
import { ApiService } from "../../../services/api.service"
import { AuthService } from "../../../services/auth/auth.service"
import { useState } from "react"
import StatusMessage from "../../../components/StatusMessage"
import useCountdown from "../../../hooks/useCountdown"

// Form Schema
const schema = z.object({
    email: z.email("Invalid email"),
})

export default function VerifyEmailForm() {
    // Hooks
    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    return (
        <form className="w-[70vw] max-w-96 flex flex-col gap-5">
            <FormStringInput
                type="email"
                label="Email Address"
                required
                register={register("email")}
                warn={valErrors.email !== undefined}
                warningMsg={valErrors.email?.message}
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
    // Hooks
    const [countdown, startCountdown] = useCountdown(30)
    const [responseMsg, setResponseMsg] = useState<
        { msg: string; variant: "success" | "error" } | undefined
    >(undefined)

    // services
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

            // Prevent user from spaming requests
            startCountdown()
        } catch (err) {
            if (err instanceof Error)
                setResponseMsg({ msg: err.message, variant: "error" })
        }
    }

    return (
        <section className="w-full flex flex-col items-center">
            <StatusMessage
                msg={responseMsg?.msg}
                variant={responseMsg?.variant ?? "info"}
                className="text-2xl"
            />
            <MyButton
                variant="positive"
                type="submit"
                disabled={hasValidationErr || isSubmitting || countdown !== 0}
                onClick={handleSubmit(submitHandler)}
                className="w-full md:w-44"
            >
                Send{countdown !== 0 && ` (${countdown})`}
            </MyButton>
        </section>
    )
}
