import React from "react"
import {Button} from "@/components/ui/button"
import {ChevronsLeft, ChevronsRight} from "lucide-react"

interface DateSelectorProps {
    year: number
    month: number
    onChange: (year: number, month: number) => void
}

const DateSelector: React.FC<DateSelectorProps> = ({year, month, onChange}) => {
    const min = 2021
    const max = 2026

    const months = (y: number): number[] => {
        if (y === 2021) return Array.from({length: 10}, (_, i) => i + 3)
        if (y === 2026) return [1, 2]
        return Array.from({length: 12}, (_, i) => i + 1)
    }

    const mList = months(year)

    const prev = () => {
        if (year > min) {
            const y = year - 1
            const avail = months(y)
            const m = avail.includes(month) ? month : avail[0]
            onChange(y, m)
        }
    }

    const next = () => {
        if (year < max) {
            const y = year + 1
            const avail = months(y)
            const m = avail.includes(month) ? month : avail[0]
            onChange(y, m)
        }
    }

    return (
        <div className="flex flex-col items-center gap-2 mb-4">
            {/* 年切り替えの部分 */}
            <div className="flex items-center gap-2">
                <Button variant="ghost" onClick={prev} disabled={year <= min}>
                    <ChevronsLeft className="w-5 h-5"/>
                </Button>
                <span className="text-lg font-semibold text-gray-800">{year}年</span>
                <Button variant="ghost" onClick={next} disabled={year >= max}>
                    <ChevronsRight className="w-5 h-5"/>
                </Button>
            </div>

            {/* 月切り替えの部分 */}
            <div className="flex flex-row items-center justify-center gap-3 flex-wrap">
                {mList.map((m) => (
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
