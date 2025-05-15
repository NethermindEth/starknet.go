import { LineChartIllustration } from "../../public/images/LineChartIllustration"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeaderCell,
  TableRoot,
  TableRow,
} from "../Table"

const summary = [
  {
    name: "North Field Corn",
    value: "21,349 bu",
    planted: "19,000 bu",
    water: "14,033 gal",
    yield: "+11.2%",
    efficiency: "+7.8%",
    nutrients: "+4.9%",
    bgColor: "bg-amber-500",
    changeType: "positive",
  },
  {
    name: "West Field Soybeans",
    value: "25,943 bu",
    planted: "23,600 bu",
    water: "11,033 gal",
    yield: "+3.1%",
    efficiency: "+5.6%",
    nutrients: "+2.9%",
    bgColor: "bg-emerald-500",
    changeType: "positive",
  },
  {
    name: "South Field Wheat",
    value: "9,443 bu",
    planted: "14,600 bu",
    water: "2,033 gal",
    yield: "-5.1%",
    efficiency: "-6.3%",
    nutrients: "-9.9%",
    bgColor: "bg-yellow-400",
    changeType: "negative",
  },
]

export default function FieldPerformance() {
  return (
    <div className="h-150 shrink-0 overflow-hidden [mask-image:radial-gradient(white_30%,transparent_90%)] perspective-[4000px] perspective-origin-center">
      <div className="-translate-y-10 -translate-z-10 rotate-x-10 rotate-y-20 -rotate-z-10 transform-3d">
        <h3 className="text-sm text-gray-500">Field Yield Performance</h3>
        <p className="mt-1 text-3xl font-semibold text-gray-900">
          32,227 bushels
        </p>
        <p className="mt-1 text-sm font-medium">
          <span className="text-emerald-700">+430 bushels (4.1%)</span>{" "}
          <span className="font-normal text-gray-500">Past growing season</span>
        </p>
        <LineChartIllustration className="mt-8 w-full min-w-200 shrink-0" />

        <TableRoot className="mt-6 min-w-200">
          <Table>
            <TableHead>
              <TableRow>
                <TableHeaderCell>Field</TableHeaderCell>
                <TableHeaderCell className="text-right">Yield</TableHeaderCell>
                <TableHeaderCell className="text-right">
                  Expected
                </TableHeaderCell>
                <TableHeaderCell className="text-right">
                  Water Used
                </TableHeaderCell>
                <TableHeaderCell className="text-right">
                  Yield Diff
                </TableHeaderCell>
                <TableHeaderCell className="text-right">
                  Efficiency
                </TableHeaderCell>
                <TableHeaderCell className="text-right">
                  Nutrients
                </TableHeaderCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {summary.map((item) => (
                <TableRow key={item.name}>
                  <TableCell className="font-medium text-gray-900">
                    <div className="flex space-x-3">
                      <span
                        className={item.bgColor + " w-1 shrink-0 rounded"}
                        aria-hidden="true"
                      />
                      <span>{item.name}</span>
                    </div>
                  </TableCell>
                  <TableCell className="text-right">{item.value}</TableCell>
                  <TableCell className="text-right">{item.planted}</TableCell>
                  <TableCell className="text-right">{item.water}</TableCell>
                  <TableCell className="text-right">
                    <span
                      className={
                        item.changeType === "positive"
                          ? "text-emerald-700"
                          : "text-red-700"
                      }
                    >
                      {item.yield}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <span
                      className={
                        item.changeType === "positive"
                          ? "text-emerald-700"
                          : "text-red-700"
                      }
                    >
                      {item.efficiency}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <span
                      className={
                        item.changeType === "positive"
                          ? "text-emerald-700"
                          : "text-red-700"
                      }
                    >
                      {item.nutrients}
                    </span>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableRoot>
      </div>
    </div>
  )
}
