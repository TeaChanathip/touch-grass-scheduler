"use client"

import { useForm, UseFormHandleSubmit } from "react-hook-form"
import { selectUser } from "../../store/features/user/userSlice"
import { useAppSelector } from "../../store/hooks"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { User, UserGender } from "../../interfaces/User.interface"
import FormStringInput from "../../components/FormStringInput"
import FormSelect from "../../components/FormSelect"
import { genderOptions } from "../../constants/options"
import { CountryCodeSchema } from "../../schemas/CountryCodeSchema"
import parsePhoneNumberFromString, { PhoneNumber } from "libphonenumber-js"
import FormPhone from "../../components/FormPhone"
import MyButton from "../../components/MyButton"
import StatusMessage from "../../components/StatusMessage"
import ImageUploader from "../../components/FormImage"

export default function ProfilePage() {
    return (
        <>
            <ImageUploader
                fallBackSrc="default_avartar.svg"
                alt="avartar"
                className="bg-red-500 size-[200px] rounded-full"
            />
            <UserProfileForm />
        </>
    )
}

// Form Schema
const schema = z.object({
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
    avartar_url: z.url().optional(),
})

function UserProfileForm() {
    // Store
    const user = useAppSelector(selectUser) as User

    // Parse phone number
    const phoneNumber = parsePhoneNumberFromString(user.phone) as PhoneNumber

    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting },
        reset,
    } = useForm({
        resolver: zodResolver(schema),
        mode: "onChange",
        defaultValues: {
            first_name: user.first_name,
            middle_name: user.middle_name,
            last_name: user.last_name,
            gender: user.gender,
            country_code: phoneNumber.country,
            phone: phoneNumber.nationalNumber,
            avartar_url: user.avatar_url,
        },
    })

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
                    valErrors.country_code?.message ?? valErrors.phone?.message
                }
            />
            <ButtonSection
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
                hasValidationErr={Object.keys(valErrors).length !== 0}
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
    const submitHandler = () => {}

    return (
        <section className="flex flex-col items-center">
            <StatusMessage msg={"test"} variant="error" className="text-2xl" />
            <div className="w-full mt-3 flex flex-col md:flex-row gap-5">
                <MyButton
                    variant="positive"
                    type="submit"
                    disabled={hasValidationErr || isSubmitting}
                    onClick={handleSubmit(submitHandler)}
                    className="w-full md:w-44"
                >
                    Save
                </MyButton>
                <MyButton
                    variant="negative"
                    type="button"
                    className="w-full md:w-44"
                >
                    Reset
                </MyButton>
            </div>
        </section>
    )
}
