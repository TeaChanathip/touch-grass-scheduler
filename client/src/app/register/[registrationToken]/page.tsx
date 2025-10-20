"use client"

import * as z from "zod"
import { UserGender, UserRole } from "../../../interfaces/User.interface"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import FormStringInput from "../../../components/FormInput"
import FormSelect from "../../../components/FormSelect"
import FormRadioGroup from "../../../components/FormRadioGroup"
import { useEffect, useState } from "react"
import MyButton from "../../../components/MyButton"
import FormPhone from "../../../components/FormPhone"
import { useAppDispatch, useAppSelector } from "../../../store/hooks"
import {
    selectUserErrMsg,
    selectUserStatus,
} from "../../../store/features/user/userSlice"
import { useParams, useRouter } from "next/navigation"
import { CircularProgress } from "@mui/material"
import { register as registerAction } from "../../../store/features/user/userSlice"
import { RegisterPayload } from "../../../interfaces/RegisterPayload.interface"
import parsePhoneNumberFromString, {
    isValidPhoneNumber,
} from "libphonenumber-js"
import { CountryCodeSchema } from "../../../schemas/CountryCodeSchema"
import { genderOptions, roleOptions } from "../../../constants/options"

export default function RegisterPage() {
    // Store
    const dispatch = useAppDispatch()
    const userStatus = useAppSelector(selectUserStatus)
    const userErrMsg = useAppSelector(selectUserErrMsg)

    // Other hooks
    const [errMsg, setErrMsg] = useState("")
    const router = useRouter()
    const params = useParams<{ registrationToken: string }>()

    const schema = z
        .object({
            first_name: z
                .string()
                .nonempty("Required")
                .max(128, "Too long")
                .regex(/^[a-zA-Z]+$/, "Letters only"),
            middle_name: z
                .string()
                .regex(/^[a-zA-Z]+$/, "Letters only")
                .max(128, "Too long")
                .or(z.literal("")),
            last_name: z
                .string()
                .regex(/^[a-zA-Z]+$/, "Letters only")
                .max(128, "Too long")
                .or(z.literal("")),
            gender: z.enum(UserGender, {
                error: () => ({ message: "Select your gender" }),
            }),
            country_code: CountryCodeSchema,
            phone: z.string().regex(/^\d{7,15}$/),
            role: z.enum(UserRole, {
                error: () => ({ message: "Select your role" }),
            }),
            school_num: z
                .string()
                .regex(/^[0-9]+$/, "Numbers only")
                .min(1, "Required")
                .max(16, "Too long")
                .optional(),
            // email: z.email("Invalid email"),
            password: z.string().min(8, "Min 8 characters").max(64, "Too long"),
            confirm_password: z
                .string()
                .min(8, "Min 8 characters")
                .max(64, "Too long"),
        })
        .refine((data) => data.password === data.confirm_password, {
            message: "Passwords don't match",
            path: ["confirm_password"],
        })
        .refine((data) => isValidPhoneNumber(data.phone, data.country_code), {
            message: "Invalid phone number",
            path: ["phone"],
        })

    // Use react-hook-form
    const {
        register,
        handleSubmit,
        formState: { errors: valErrors },
        watch,
        setValue,
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    // Watch role in real-time
    const currentRole = watch("role")

    // Clear school_num if role is not school personnel
    useEffect(() => {
        if (![UserRole.STUDENT, UserRole.TEACHER].includes(currentRole)) {
            setValue("school_num", undefined)
        }
    }, [setValue, currentRole])

    // Submit Handler
    const onSubmit = (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)
        if (result.success) {
            const registerPayload: RegisterPayload = {
                role: result.data.role,
                first_name: result.data.first_name,
                middle_name: result.data.middle_name,
                last_name: result.data.last_name,
                phone: parsePhoneNumberFromString(
                    result.data.phone,
                    result.data.country_code
                )!.format("E.164"),
                gender: result.data.gender,
                password: result.data.password,
                school_num: result.data.school_num,
            }

            dispatch(
                registerAction({
                    registrationToken: params.registrationToken,
                    registerPayload: registerPayload,
                })
            )
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
            <h1 className="text-5xl">REGISTER</h1>

            <form
                onSubmit={handleSubmit(onSubmit)}
                className="w-4/5 lg:w-1/3 flex flex-col gap-3 mb-5"
            >
                <span className="flex flex-row gap-4">
                    <FormStringInput
                        label="Frist Name"
                        type="text"
                        required
                        register={register("first_name")}
                        warn={valErrors.first_name !== undefined}
                        warningMsg={valErrors.first_name?.message ?? ""}
                    />
                    <FormStringInput
                        label="Middle Name"
                        type="text"
                        register={register("middle_name")}
                        warn={valErrors.middle_name !== undefined}
                        warningMsg={valErrors.middle_name?.message ?? ""}
                    />
                </span>
                <span className="flex flex-row gap-4 justify-between">
                    <FormStringInput
                        label="Last Name"
                        type="text"
                        register={register("last_name")}
                        warn={valErrors.last_name !== undefined}
                        warningMsg={valErrors.last_name?.message ?? ""}
                    />
                    <FormSelect
                        label="Gender"
                        optionItems={genderOptions}
                        required
                        register={register("gender")}
                        warn={valErrors.gender !== undefined}
                        warningMsg={valErrors.gender?.message ?? ""}
                    />
                </span>
                <FormPhone
                    required
                    countryCodeRegister={register("country_code")}
                    phoneRegister={register("phone")}
                    warn={
                        valErrors.country_code !== undefined ||
                        valErrors.phone !== undefined
                    }
                    warningMsg={
                        valErrors.country_code?.message ??
                        valErrors.phone?.message
                    }
                />
                <FormRadioGroup
                    label="Role"
                    options={roleOptions}
                    register={register("role")}
                    warn={valErrors.role !== undefined}
                    warningMsg={valErrors.role?.message ?? ""}
                />
                {[UserRole.STUDENT, UserRole.TEACHER].includes(currentRole) && (
                    <FormStringInput
                        label="School Number"
                        type="text"
                        required
                        register={register("school_num")}
                        warn={valErrors.school_num !== undefined}
                        warningMsg={valErrors.school_num?.message ?? ""}
                    />
                )}
                {/* <FormStringInput */}
                {/*     label="Email" */}
                {/*     type="email" */}
                {/*     required */}
                {/*     register={register("email")} */}
                {/*     warn={valErrors.email !== undefined} */}
                {/*     warningMsg={valErrors.email?.message ?? ""} */}
                {/* /> */}
                <FormStringInput
                    label="Password"
                    type="password"
                    required
                    register={register("password")}
                    warn={valErrors.password !== undefined}
                    warningMsg={valErrors.password?.message ?? ""}
                />
                <FormStringInput
                    label="Confirm Password"
                    type="password"
                    required
                    register={register("confirm_password")}
                    warn={valErrors.confirm_password !== undefined}
                    warningMsg={valErrors.confirm_password?.message ?? ""}
                />
                <MyButton
                    variant="positive"
                    type="submit"
                    disabled={Object.keys(valErrors).length > 0}
                    className="mt-5 w-full md:w-44 self-center"
                >
                    Register
                </MyButton>
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

// TODO: Send verification email
// TODO: Refactor API to not expose that the email already exists in client.
// Instead, sending a warning email.
