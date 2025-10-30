"use client"

import {
    useForm,
    UseFormHandleSubmit,
    UseFormRegister,
    UseFormReset,
} from "react-hook-form"
import {
    selectUser,
    selectUserErrMsg,
    userUpdateProfile,
} from "../../store/features/user/userSlice"
import { useAppDispatch, useAppSelector } from "../../store/hooks"
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
import ImageUploader from "../../components/ImageUploader"
import { Dispatch, SetStateAction, useState } from "react"
import { Visibility } from "@mui/icons-material"
import { idnEmail } from "zod/v4/core/regexes.cjs"
import { isSchoolPersonnel } from "../../utils/isSchoolPersonnel"

export default function ProfilePage() {
    // Hooks
    const [isEditing, setEditing] = useState(false)

    return (
        <>
            <ImageUploader
                fallBackSrc="default_avartar.svg"
                alt="avartar"
                disabled={!isEditing}
                width={200}
                height={200}
                className="rounded-full border border-[1px] border-prim-green-800 bg-prim-green-50"
            />
            <UserProfileForm isEditing={isEditing} setEditing={setEditing} />
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

    // Fixed fields
    email: z.email(),
    role: z.string(),
    school_num: z.string().optional(),
})

function UserProfileForm({
    isEditing,
    setEditing,
}: {
    isEditing: boolean
    setEditing: Dispatch<SetStateAction<boolean>>
}) {
    // Store
    const user = useAppSelector(selectUser) as User

    // Parse phone number
    const phoneNumber = parsePhoneNumberFromString(user.phone) as PhoneNumber

    const {
        register,
        handleSubmit,
        formState: { errors: valErrors, isSubmitting, isDirty, dirtyFields },
        reset,
    } = useForm({
        resolver: zodResolver(schema),
        mode: "onChange",
        defaultValues: {
            first_name: user.first_name,
            middle_name: user.middle_name ?? "",
            last_name: user.last_name ?? "",
            gender: user.gender,
            country_code: phoneNumber.country,
            phone: phoneNumber.nationalNumber,
            email: user.email,
            role: user.role,
            school_num: user.school_num,
        },
    })

    return (
        <form className="w-4/5 lg:w-1/3 flex flex-col gap-3 mb-5">
            <div className="w-full flex flex-row gap-4 justify-between">
                <FormStringInput
                    label="Frist Name"
                    type="text"
                    required
                    register={register("first_name")}
                    warn={valErrors.first_name !== undefined}
                    warningMsg={valErrors.first_name?.message}
                    readOnly={!isEditing}
                    hideMsg={!isEditing}
                />
                <FormStringInput
                    label="Middle Name"
                    type="text"
                    register={register("middle_name")}
                    warn={valErrors.middle_name !== undefined}
                    warningMsg={valErrors.middle_name?.message}
                    readOnly={!isEditing}
                    hideMsg={!isEditing}
                />
            </div>
            <div className="w-full flex flex-row gap-4 justify-between">
                <FormStringInput
                    label="Last Name"
                    type="text"
                    register={register("last_name")}
                    warn={valErrors.last_name !== undefined}
                    warningMsg={valErrors.last_name?.message}
                    readOnly={!isEditing}
                    hideMsg={!isEditing}
                />
                <FormSelect
                    label="Gender"
                    optionItems={genderOptions}
                    required
                    register={register("gender")}
                    warn={valErrors.gender !== undefined}
                    warningMsg={valErrors.gender?.message ?? ""}
                    disabled={!isEditing}
                    hideMsg={!isEditing}
                />
            </div>
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
                readOnly={!isEditing}
                hideMsg={!isEditing}
            />
            <UnmodifiedFieldsGroup isEditing={isEditing} register={register} />
            <ButtonSection
                handleSubmit={handleSubmit}
                isSubmitting={isSubmitting}
                hasValidationErr={Object.keys(valErrors).length !== 0}
                isDirty={isDirty}
                dirtyFields={dirtyFields}
                reset={reset}
                isEditing={isEditing}
                setEditing={setEditing}
            />
        </form>
    )
}

function UnmodifiedFieldsGroup({
    isEditing,
    register,
}: {
    isEditing: boolean
    register: UseFormRegister<z.infer<typeof schema>>
}) {
    // Store
    const user = useAppSelector(selectUser) as User

    return (
        <div hidden={isEditing} className="flex flex-col gap-4">
            <div className="flex flex-col md:flex-row justify-between gap-4">
                <FormStringInput
                    label="Email"
                    type="text"
                    register={register("email")}
                    readOnly
                    hideMsg
                />
                <FormStringInput
                    label="Role"
                    type="text"
                    register={register("role")}
                    readOnly
                    hideMsg
                />
            </div>
            {isSchoolPersonnel(user.role) && (
                <FormStringInput
                    label="School Number"
                    type="text"
                    register={register("school_num")}
                    readOnly
                    hideMsg
                />
            )}
        </div>
    )
}

function ButtonSection({
    handleSubmit,
    isSubmitting,
    hasValidationErr,
    isDirty,
    dirtyFields,
    reset,
    isEditing,
    setEditing,
}: {
    handleSubmit: UseFormHandleSubmit<z.infer<typeof schema>>
    isSubmitting: boolean
    hasValidationErr: boolean
    isDirty: boolean
    dirtyFields: Partial<
        Readonly<{
            first_name?: boolean
            middle_name?: boolean
            last_name?: boolean
            gender?: boolean
            country_code?: boolean
            phone?: boolean
        }>
    >
    reset: UseFormReset<z.infer<typeof schema>>
    isEditing: boolean
    setEditing: Dispatch<SetStateAction<boolean>>
}) {
    // Store
    const dispatch = useAppDispatch()
    const userErrMsg = useAppSelector(selectUserErrMsg)

    // Button Handler
    const submitHandler = async (formData: z.infer<typeof schema>) => {
        const result = schema.safeParse(formData)

        if (result.error) return

        // update only dirty fields
        await dispatch(
            userUpdateProfile({
                ...(dirtyFields.first_name && {
                    first_name: result.data.first_name,
                }),
                ...(dirtyFields.middle_name && {
                    middle_name: result.data.middle_name,
                }),
                ...(dirtyFields.last_name && {
                    last_name: result.data.last_name,
                }),
                ...(dirtyFields.gender && {
                    gender: result.data.gender,
                }),
                ...((dirtyFields.country_code || dirtyFields.phone) && {
                    phone: parsePhoneNumberFromString(
                        result.data.phone,
                        result.data.country_code
                    )!.format("E.164"),
                }),
            })
        )
    }

    const dismissHandler = () => {
        reset()
        setEditing(false)
    }

    return (
        <section className="flex flex-col items-center">
            <StatusMessage
                msg={userErrMsg}
                variant="error"
                className="text-2xl bt-3"
            />
            {isEditing ? (
                <div className="w-full flex flex-col justify-center md:flex-row gap-5">
                    <MyButton
                        variant="positive"
                        type="submit"
                        disabled={!isDirty || hasValidationErr || isSubmitting}
                        onClick={handleSubmit(submitHandler)}
                        className="w-full md:w-44"
                    >
                        Save
                    </MyButton>
                    <MyButton
                        variant="negative"
                        type="button"
                        onClick={() => dismissHandler()}
                        className="w-full md:w-44"
                    >
                        Dimiss
                    </MyButton>
                </div>
            ) : (
                <MyButton
                    variant="positive"
                    type="button"
                    onClick={() => setEditing(true)}
                    className="w-full md:w-44"
                >
                    Edit
                </MyButton>
            )}
        </section>
    )
}

// TODO: Complete image uploader
// TODO: Update the logic for showing response message (must include the error from uploading an image)
