import { UseFormRegisterReturn } from "react-hook-form"
import FormStringInput from "./FormInput"
import FormSelect from "./FormSelect"

import CountryList from "country-list-with-dial-code-and-flag"

const dialCodeOptions = CountryList.getAll()
    .map((country) => {
        return {
            value: country.dial_code,
            label: `${country.flag} ${country.dial_code}`,
            id: `select-dial-code-${country.name}-${country.dial_code}`,
        }
    })
    .sort(({ value: value1 }, { value: value2 }) => {
        return Number(value1.slice(1)) - Number(value2.slice(1))
    })

export default function FormPhone({
    required,
    dialCodeRegister,
    phoneRegister,
    warn,
    warningMsg,
}: {
    required?: boolean
    dialCodeRegister?: UseFormRegisterReturn<any>
    phoneRegister?: UseFormRegisterReturn<any>
    warn?: boolean
    warningMsg?: string
}) {
    return (
        <div className="flex flex-col">
            <label className="text-2xl">
                Phone{required && <span className="text-prim-red">*</span>}
            </label>
            <span className="flex flex-row gap-4 w-full">
                <span className="w-48">
                    <FormSelect
                        optionItems={dialCodeOptions}
                        required
                        register={dialCodeRegister}
                        warn={warn}
                    />
                </span>
                <span className="w-full">
                    <FormStringInput
                        type="tel"
                        required
                        register={phoneRegister}
                        warn={warn}
                    />
                </span>
            </span>
            {warningMsg !== undefined && (
                <p className="self-center text-prim-red">{warningMsg}&nbsp;</p>
            )}
        </div>
    )
}
