import { UseFormRegisterReturn } from "react-hook-form"
import FormStringInput from "./FormStringInput"
import { countryCodesOptions } from "../constants/options"
import FormDatalist from "./FormDatalist"
import StatusMessage from "./StatusMessage"

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
                        hideMsg
                    />
                </span>
                <span className="w-full">
                    <FormStringInput
                        type="tel"
                        required
                        register={phoneRegister}
                        warn={warn}
                        hideMsg
                    />
                </span>
            </span>
            <StatusMessage
                msg={warningMsg}
                variant="error"
                className="text-xl"
            />
        </div>
    )
}
