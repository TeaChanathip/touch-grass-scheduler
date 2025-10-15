"use client"

import * as z from "zod"
import { UserGender, UserRole } from "../../interfaces/User.interface"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import FormStringInput from "../../components/FormInput"
import FormSelect from "../../components/FormSelect"
import snakeToTitleCase from "../../utils/snakeToTitleCase"
import FormRadioGroup from "../../components/FormRadioGroup"
import { useEffect } from "react"
import MyButton from "../../components/MyButton"
import FormPhone from "../../components/FormPhone"

// Generate options from Enum
const genderOptions = Object.keys(UserGender).map((key) => {
    return {
        value: UserGender[key],
        label: snakeToTitleCase(UserGender[key]),
        id: `select-gender-${UserGender[key]}`,
    }
})

const roleOptions = Object.keys(UserRole)
    .filter((key) => UserRole[key] !== "admin")
    .map((key) => {
        return {
            value: UserRole[key],
            label: snakeToTitleCase(UserRole[key]),
            id: `select-role-${UserRole[key]}`,
        }
    })

export default function Register() {
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
            dial_code: z.string().regex(/^\+\d{1,4}$/),
            phone: z.string().regex(/^\d{7,14}$/),
            role: z.enum(UserRole, {
                error: () => ({ message: "Select your role" }),
            }),
            school_num: z
                .string()
                .regex(/^[0-9]+$/, "Numbers only")
                .min(1, "Required")
                .max(16, "Too long")
                .optional(),
            email: z.email("Invalid email"),
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
        console.log("test")
    }

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
                {/* <FormStringInput
                    label="Phone"
                    type="tel"
                    required
                    register={register("phone")}
                    warningMsg={valErrors.phone?.message}
                /> */}
                <FormPhone
                    required
                    dialCodeRegister={register("dial_code")}
                    phoneRegister={register("phone")}
                    warn={
                        valErrors.dial_code !== undefined ||
                        valErrors.phone !== undefined
                    }
                    warningMsg={
                        valErrors.dial_code || valErrors.phone
                            ? "Invalid phone format"
                            : ""
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
                <FormStringInput
                    label="Email"
                    type="email"
                    required
                    register={register("email")}
                    warn={valErrors.email !== undefined}
                    warningMsg={valErrors.email?.message ?? ""}
                />
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
        </div>
    )
}
