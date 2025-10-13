"use client"

import { useForm } from "react-hook-form"
import FormStringInput from "../../components/FormInput"
import MyButton from "../../components/MyButton"
import Link from "next/link"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

export default function Login() {
    const schema = z.object({
        email: z.email("Invalid email format"),
        password: z
            .string()
            .min(8, "At least 8 characters")
            .max(64, "At most 64 characters"),
    })

    const {
        register,
        handleSubmit,
        formState: { errors: valErrors },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    const onSubmit = (formData: z.infer<typeof schema>) => {}

    return (
        <div className="flex flex-col gap-10 items-center pt-10">
            <h1 className="text-5xl">LOGIN</h1>

            <form
                onSubmit={handleSubmit(onSubmit)}
                className="w-[70vw] max-w-96 flex flex-col gap-5"
            >
                <FormStringInput
                    label="Email Address"
                    type="email"
                    required
                    register={register("email")}
                    warningMsg={valErrors.email?.message}
                />
                <FormStringInput
                    label="Password"
                    type="password"
                    required
                    register={register("password")}
                    warningMsg={valErrors.password?.message}
                />
                <Link
                    href="/forgot-password"
                    className="w-fit self-center underline"
                >
                    Forgot the password?
                </Link>
                <div className="mt-3 flex flex-col md:flex-row gap-5 justify-center">
                    <MyButton
                        variant="positive"
                        type="submit"
                        className="w-full md:w-44"
                        disabled={Object.keys(valErrors).length > 0}
                    >
                        Login
                    </MyButton>
                    <MyButton
                        variant="neutral"
                        type="button"
                        className="w-full md:w-44"
                    >
                        Register
                    </MyButton>
                </div>
            </form>
        </div>
    )
}

// TODO: disable login button when the form is invalid
