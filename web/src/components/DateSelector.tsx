import React from "react"
import {Button} from "@/components/ui/button"
import {ChevronsLeft, ChevronsRight} from "lucide-react"

interface DateSelectorProps {
    year: number
    month: number
    onChange: (year: number, month: number) => void
}

const DateSelector: React.FC<DateSelectorProps> = ({year, month, onChange}) => {
    const handlePrevYear = () => {
        onChange(year - 1, month)
    }
    const handleNextYear = () => {
        onChange(year + 1, month)
    }
    const months = Array.from({length: 12}, (_, i) => i + 1)

    return (
        <div className="flex flex-col items-center gap-2 mb-4">
            {/* 年切り替えの部分 */}
            <div className="flex items-center gap-2">
                <Button variant="ghost" onClick={handlePrevYear}>
                    <ChevronsLeft className="w-5 h-5"/>
                </Button>
                <span className="text-lg font-semibold text-gray-800">{year}年</span>
                <Button variant="ghost" onClick={handleNextYear}>
                    <ChevronsRight className="w-5 h-5"/>
                </Button>
            </div>

            {/* 月切り替えの部分 */}
            <div className="flex flex-row items-center justify-center gap-3">
                {months.map((m) => (
                    <Button
                        key={m}
                        variant={m === month ? "default" : "ghost"}
                        onClick={() => onChange(year, m)}
                    >
                        {m}月
                    </Button>
                ))}
            </div>
        </div>
    )
}

export default DateSelector
