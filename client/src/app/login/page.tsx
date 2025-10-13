"use client"

import FormStringInput from "../../components/FormInput"

import * as z from "zod"
import MyButton from "../../components/MyButton"
import Link from "next/link"

export default function Login() {
    const schema = z.object({
        email: z.email("Invalid email format"),
        password: z
            .string()
            .min(8, "Must be at least 8 characters")
            .max(64, "Must not longer than 64 characters"),
    })

    const formAction = (formData: FormData) => {}

    return (
        <div className="flex flex-col gap-10 items-center pt-10">
            <h1 className="text-5xl">LOGIN</h1>

            <form
                action={formAction}
                className="w-[70vw] max-w-96 flex flex-col gap-5"
            >
                <FormStringInput
                    label="Email Address"
                    name="email"
                    type="email"
                    required
                    schema={schema.shape.email}
                />
                <FormStringInput
                    label="Password"
                    name="password"
                    type="password"
                    required
                    schema={schema.shape.password}
                />
                <Link
                    href="/forgot-password"
                    className="w-fit self-center underline"
                >
                    Forgot the password?
                </Link>
                <div className="mt-6 flex flex-col md:flex-row gap-5 justify-center">
                    <MyButton
                        variant="positive"
                        type="submit"
                        className="w-full md:w-44"
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
