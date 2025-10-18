"use client"

import { useForm } from "react-hook-form"
import FormStringInput from "../../components/FormInput"
import MyButton from "../../components/MyButton"
import Link from "next/link"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { useAppDispatch, useAppSelector } from "../../store/hooks"
import {
    login,
    selectUserErrMsg,
    selectUserStatus,
} from "../../store/features/user/userSlice"
import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { CircularProgress } from "@mui/material"

export default function LoginPage() {
    // Store
    const dispatch = useAppDispatch()
    const userStatus = useAppSelector(selectUserStatus)
    const userErrMsg = useAppSelector(selectUserErrMsg)

    const [errMsg, setErrMsg] = useState("")
    const router = useRouter()

    // Define schema for validation
    const schema = z.object({
        email: z.email("Invalid email format"),
        password: z
            .string()
            .min(8, "At least 8 characters")
            .max(64, "At most 64 characters"),
    })

    // Use react-hook-form
    const {
        register,
        handleSubmit,
        formState: { errors: valErrors },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    // Submit Handler
    const onSubmit = (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)
        if (result.success) {
            dispatch(login(result.data))
        } else {
            setErrMsg("Validation Error")
        }
    }

    useEffect(() => {
        switch (userStatus) {
            case "authenticated":
                router.push("/")
                break
            case "unauthenticated":
                setErrMsg("Incorrect Email or Password")
                break
            case "error":
                setErrMsg(userErrMsg ?? "Unknown Error")
                break
            default:
                setErrMsg("")
                break
        }
    }, [userStatus, userErrMsg, router])

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
                    warn={valErrors.email != undefined}
                    warningMsg={valErrors.email?.message ?? ""}
                />
                <FormStringInput
                    label="Password"
                    type="password"
                    required
                    register={register("password")}
                    warn={valErrors.password != undefined}
                    warningMsg={valErrors.password?.message ?? ""}
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
                        disabled={
                            Object.keys(valErrors).length > 0 ||
                            userStatus == "loading"
                        }
                    >
                        Login
                    </MyButton>
                    <MyButton
                        variant="neutral"
                        type="button"
                        className="w-full md:w-44"
                        onClick={() => router.push("/register")}
                    >
                        Register
                    </MyButton>
                </div>
            </form>
            {userStatus != "loading" && errMsg && (
                <p className="text-prim-red">{errMsg}</p>
            )}
            {userStatus == "loading" && (
                <span className="text-prim-green-400">
                    <CircularProgress color="inherit" />
                </span>
            )}
        </div>
    )
}
