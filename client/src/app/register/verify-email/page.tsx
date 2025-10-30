"use client"

import { FormProvider, useForm, useFormContext } from "react-hook-form"
import FormStringInput from "../../../components/FormStringInput"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import MyButton from "../../../components/MyButton"
import { ApiService } from "../../../services/api.service"
import { AuthService } from "../../../services/auth/auth.service"
import { memo, useState } from "react"
import StatusMessage from "../../../components/StatusMessage"
import useCountdown from "../../../hooks/useCountdown"

// Form Schema
const schema = z.object({
    email: z.email("Invalid email"),
})

export default function VerifyEmailForm() {
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
                <ButtonSection />
            </form>
        </FormProvider>
    )
}

const ButtonSection = memo(function ButtonSection() {
    // Form Context
    const {
        handleSubmit,
        formState: { isValid, isSubmitting },
    } = useFormContext()

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
                disabled={!isValid || isSubmitting || countdown !== 0}
                onClick={handleSubmit(submitHandler)}
                className="w-full md:w-44"
            >
                Send{countdown !== 0 && ` (${countdown})`}
            </MyButton>
        </section>
    )
})
