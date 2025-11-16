import axios from "axios";
import { useEffect, useState } from "react";
import type { ClickStats, Tabs } from "../../types";
import { toast } from "sonner";
import { formatError, linkShortIdToUrl } from "../../utils";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from "react-chartjs-2";
import { chartBackgroundColors, chartBorderColors } from "../../constants";

type Props = {
    shortLinkId: string;
    setTabFn: (tabName: Tabs) => void;
};

ChartJS.register(ArcElement, Tooltip, Legend);

export default function StatsTab({ shortLinkId, setTabFn }: Readonly<Props>) {
    const shortUrl = linkShortIdToUrl(shortLinkId);
    const [clickStats, setClickStats] = useState<ClickStats>([]);
    const [isLoading, setIsLoading] = useState<boolean>(false);

    useEffect(() => {
        (async () => {
            try {
                setIsLoading(true);
                const response = await axios.get<ClickStats>(`/api/stats/${shortLinkId}`);
                setClickStats(response.data);
            } catch (error) {
                console.error("Error fetching stats", error);
                toast.error(`Error fetching stats${formatError(error)}`);
            } finally {
                setIsLoading(false);
            }
        })();
    }, [shortLinkId]);

    return (
        <>
            <h2 className="text-xl font-semibold">{shortUrl}</h2>
            <p>Visits this week by referrers:</p>
            {isLoading ? (
                "Loading..."
            ) : (
                <div className="max-w-full">
                    <Pie
                        data={{
                            labels: clickStats.map((item) => item.referrer || "No referrer"),
                            datasets: [
                                {
                                    label: "Visits this week",
                                    data: clickStats.map((item) => item.count),
                                    backgroundColor: chartBackgroundColors,
                                    borderColor: chartBorderColors
                                }
                            ]
                        }}
                    />
                </div>
            )}
            <button
                className="px-1.5 py-1 bg-accent rounded-sm shadow-sm transition-all hover:brightness-110 cursor-pointer mt-2"
                onClick={() => setTabFn("shorten")}
            >
                Shorten another URL
            </button>
        </>
    );
}
