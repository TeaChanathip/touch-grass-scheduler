"use client"

import * as z from "zod"
import { UserGender, UserRole } from "../../../interfaces/User.interface"
import { useForm, UseFormHandleSubmit } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import FormStringInput from "../../../components/FormStringInput"
import FormSelect from "../../../components/FormSelect"
import FormRadioGroup from "../../../components/FormRadioGroup"
import { useEffect } from "react"
import MyButton from "../../../components/MyButton"
import FormPhone from "../../../components/FormPhone"
import { useAppDispatch, useAppSelector } from "../../../store/hooks"
import {
    selectUserErrMsg,
    selectUserStatus,
} from "../../../store/features/user/userSlice"
import { useParams, useRouter } from "next/navigation"
import { userRegister } from "../../../store/features/user/userSlice"
import { RegisterPayload } from "../../../interfaces/Auth.interface"
import parsePhoneNumberFromString, {
    isValidPhoneNumber,
} from "libphonenumber-js"
import { CountryCodeSchema } from "../../../schemas/CountryCodeSchema"
import { genderOptions, roleOptions } from "../../../constants/options"
import StatusMessage from "../../../components/StatusMessage"

export default function RegisterPage() {
    // Store
    const userStatus = useAppSelector(selectUserStatus)

    // Hooks
    const router = useRouter()

    useEffect(() => {
        if (userStatus === "authenticated") {
            router.push("/")
        }
    }, [userStatus, router])

    return <RegisterForm />
}

// Form Schema
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
        password: z.string().min(8, "Min 8 characters").max(64, "Too long"),
        confirm_password: z.string(),
    })
    .refine((data) => data.password === data.confirm_password, {
        message: "Passwords don't match",
        path: ["confirm_password"],
    })
    .refine((data) => isValidPhoneNumber(data.phone, data.country_code), {
        message: "Invalid phone number",
        path: ["phone"],
    })

function RegisterForm() {
    // Hooks
    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting },
        watch,
        setValue,
        trigger,
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    // Watch
    const currentRole = watch("role")
    const currentPwd = watch("password")

    // Clear school_num if role is not school personnel
    useEffect(() => {
        if (![UserRole.STUDENT, UserRole.TEACHER].includes(currentRole)) {
            setValue("school_num", undefined)
        }
        trigger("school_num")
    }, [setValue, currentRole])

    // Trigger confirm_password validation when password is changed
    useEffect(() => {
        if (currentPwd) {
            trigger("confirm_password")
        }
    }, [currentPwd, trigger])

    // Constants
    const schoolPersonnelRoles: UserRole[] = [
        UserRole.STUDENT,
        UserRole.TEACHER,
    ]

    return (
        <form className="w-4/5 lg:w-1/3 flex flex-col gap-3 mb-5">
            <span className="w-full flex flex-row gap-4 justify-between">
                <FormStringInput
                    label="Frist Name"
                    type="text"
                    required
                    register={register("first_name")}
                    warn={valErrors.first_name !== undefined}
                    warningMsg={valErrors.first_name?.message}
                />
                <FormStringInput
                    label="Middle Name"
                    type="text"
                    register={register("middle_name")}
                    warn={valErrors.middle_name !== undefined}
                    warningMsg={valErrors.middle_name?.message}
                />
            </span>
            <span className="w-full flex flex-row gap-4 justify-between">
                <FormStringInput
                    label="Last Name"
                    type="text"
                    register={register("last_name")}
                    warn={valErrors.last_name !== undefined}
                    warningMsg={valErrors.last_name?.message}
                />
                <FormSelect
                    label="Gender"
                    optionItems={genderOptions}
                    required
                    register={register("gender")}
                    warn={valErrors.gender !== undefined}
                    warningMsg={valErrors.gender?.message}
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
                    valErrors.country_code?.message ?? valErrors.phone?.message
                }
            />
            <FormRadioGroup
                label="Role"
                options={roleOptions}
                register={register("role")}
                warn={valErrors.role !== undefined}
                warningMsg={valErrors.role?.message}
            />
            {schoolPersonnelRoles.includes(currentRole) && (
                <FormStringInput
                    label="School Number"
                    type="text"
                    required
                    register={register("school_num")}
                    warn={valErrors.school_num !== undefined}
                    warningMsg={valErrors.school_num?.message}
                />
            )}
            <FormStringInput
                label="Password"
                type="password"
                required
                register={register("password")}
                warn={valErrors.password !== undefined}
                warningMsg={valErrors.password?.message}
            />
            <FormStringInput
                label="Confirm Password"
                type="password"
                required
                register={register("confirm_password")}
                warn={valErrors.confirm_password !== undefined}
                warningMsg={valErrors.confirm_password?.message}
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
    // Store
    const dispatch = useAppDispatch()
    const userErrMsg = useAppSelector(selectUserErrMsg)

    // Hooks
    const params = useParams<{ registrationToken: string }>()

    const submitHandler = (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)

        if (result.error) return // Error should not be possible here

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
            ...(result.data.school_num && {
                school_num: result.data.school_num,
            }),
        }

        dispatch(
            userRegister({
                registrationToken: params.registrationToken,
                registerPayload: registerPayload,
            })
        )
    }

    return (
        <section className="mt-3 flex flex-col items-center">
            <StatusMessage
                msg={userErrMsg}
                variant="error"
                className="text-2xl"
            />
            <MyButton
                variant="positive"
                type="submit"
                disabled={isSubmitting || hasValidationErr}
                onClick={handleSubmit(submitHandler)}
                className="mt-3 w-full md:w-44 self-center"
            >
                Register
            </MyButton>
        </section>
    )
}
