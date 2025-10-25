"use client"

import { useForm } from "react-hook-form"
import { selectUser } from "../../store/features/user/userSlice"
import { useAppDispatch, useAppSelector } from "../../store/hooks"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { UserGender } from "../../interfaces/User.interface"
import FormStringInput from "../../components/FormStringInput"
import FormSelect from "../../components/FormSelect"
import { genderOptions } from "../../constants/options"

export default function ProfilePage() {
    const dispatch = useAppDispatch()
    const user = useAppSelector(selectUser)!

    const schema = z.object({
        first_name: z.string().default(user.first_name),
        middle_name: z.string().default(user.middle_name),
        last_name: z.string().default(user.last_name),
        gender: z.enum(UserGender).default(user.gender),
        dial_code: z.string(),
        phone: z.string(),
        avartar_url: z.url().default(user.avatar_url),
    })

    const {
        register,
        handleSubmit,
        formState: { errors: valErrors },
    } = useForm({ resolver: zodResolver(schema), mode: "onChange" })

    return (
        <div className="flex flex-col gap-10 items-center pt-10">
            <h1 className="text-5xl">Profile</h1>
            <form className="w-4/5 lg:w-1/3 flex flex-col gap-3 mb-5">
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
            </form>
        </div>
    )
}
