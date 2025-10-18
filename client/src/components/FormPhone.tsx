import { UseFormRegisterReturn } from "react-hook-form"
import FormStringInput from "./FormInput"
import { countryCodesOptions } from "../constants/options"
import FormDatalist from "./FormDatalist"

export default function FormPhone({
    required,
    countryCodeRegister,
    phoneRegister,
    warn,
    warningMsg,
}: {
    required?: boolean
    countryCodeRegister?: UseFormRegisterReturn<any>
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
                <span className="w-20">
                    <FormDatalist
                        optionItems={countryCodesOptions}
                        required
                        register={countryCodeRegister}
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
