import { useFormContext } from "react-hook-form"
import FormStringInput from "./FormStringInput"
import { countryCodesOptions } from "../constants/options"
import FormDatalist from "./FormDatalist"
import StatusMessage from "./StatusMessage"
import { memo } from "react"

const FormPhone = ({
    countryCodeName,
    phoneName,
    readOnly,
    required,
    warn,
    hideMsg,
}: {
    countryCodeName: string
    phoneName: string
    readOnly?: boolean
    required?: boolean
    warn?: boolean
    hideMsg?: boolean
}) => {
    // Form Context
    const {
        formState: { errors },
    } = useFormContext()
    const warningMsg =
        (errors[countryCodeName]?.message as string | undefined) ||
        (errors[phoneName]?.message as string | undefined)

    return (
        <div className="flex flex-col">
            <label className="text-2xl">
                Phone{required && <span className="text-prim-red">*</span>}
            </label>
            <span className="flex flex-row gap-4 w-full">
                <span className="w-24">
                    <FormDatalist
                        optionItems={countryCodesOptions}
                        name={countryCodeName}
                        readOnly={readOnly}
                        required={required}
                        warn={warn}
                        hideMsg
                    />
                </span>
                <span className="w-full">
                    <FormStringInput
                        type="tel"
                        name={phoneName}
                        readOnly={readOnly}
                        required={required}
                        warn={warn}
                        hideMsg
                    />
                </span>
            </span>
            {!hideMsg && (
                <StatusMessage
                    msg={warningMsg}
                    variant="error"
                    className="text-xl"
                />
            )}
        </div>
    )
}

export default memo(FormPhone)
