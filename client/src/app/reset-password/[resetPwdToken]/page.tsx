"use client"

import { FormProvider, useForm, useFormContext } from "react-hook-form"
import * as z from "zod"
import FormPassword from "../../../components/FormPassword"
import { zodResolver } from "@hookform/resolvers/zod"
import { memo, useState } from "react"
import { authService } from "../../../services/auth/auth.service"
import StatusMessage from "../../../components/StatusMessage"
import MyButton from "../../../components/MyButton"
import { useParams } from "next/navigation"

// Form Schema
const schema = z
    .object({
        password: z.string().min(8, "Min 8 characters").max(64, "Too long"),
        confirm_password: z.string(),
    })
    .refine((data) => data.password === data.confirm_password, {
        message: "Passwords don't match",
        path: ["confirm_password"],
    })

export default function RestPwdPage() {
    const formMethods = useForm({
        resolver: zodResolver(schema),
        mode: "onChange",
    })

    return (
        <FormProvider {...formMethods}>
            <form className="w-[70vw] max-w-96 flex flex-col gap-5">
                <FormPassword label="Password" name="password" />
                <FormPassword
                    label="Confirm Password"
                    name="confirm_password"
                />
            </form>
            <ButtonSection />
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
    const [responseMsg, setResponseMsg] = useState<
        { msg: string; variant: "success" | "error" } | undefined
    >(undefined)
    const params = useParams<{ resetPwdToken: string }>()

    const submitHandler = async (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)

        if (result.error) return

        try {
            await authService.resetPwd({
                reset_pwd_token: params.resetPwdToken,
                new_password: result.data.password,
            })
            setResponseMsg({
                msg: "Password reset successfully",
                variant: "success",
            })
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
                disabled={!isValid || isSubmitting}
                onClick={handleSubmit(submitHandler)}
                className="w-full md:w-44"
            >
                Submit
            </MyButton>
        </section>
    )
})
