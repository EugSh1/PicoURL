export type ShortenLinkResponse = {
    shortUrlId: string;
};

export type Tabs = "shorten" | "success" | "stats";

export type ClickStats = {
    referrer: string | null;
    count: number;
}[];
