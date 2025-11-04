"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { FormProvider, useForm, useFormContext } from "react-hook-form"
import FormStringInput from "../../components/FormStringInput"
import { memo, useState } from "react"
import useCountdown from "../../hooks/useCountdown"
import * as z from "zod"
import StatusMessage from "../../components/StatusMessage"
import MyButton from "../../components/MyButton"
import { authService } from "../../services/auth/auth.service"

// Form Schema
const schema = z.object({
    email: z.email("Invalid email"),
})

export default function ForgotPwdPage() {
    const formMethods = useForm({
        resolver: zodResolver(schema),
        mode: "onChange",
    })

    return (
        <FormProvider {...formMethods}>
            <p className="text-lg text-center text-wrap w-[70vw] lg:max-w-[40vw] bg-prim-green-50 rounded-2xl py-3 px-5 drop-shadow-md">
                <span className="font-semibold">Don&apos;t worry!&nbsp;</span>
                We know that remember a password is not easy. Please enter an
                email of your account in the form below, so we can send you a
                link to reset the password.
            </p>
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

    const submitHandler = async (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)

        try {
            await authService.getResetPwdMail(result.data!.email)
            setResponseMsg({
                msg: "Reset password email was sent",
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
