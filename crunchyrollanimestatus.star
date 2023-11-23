"""
Applet: Crunchyroll Anime Status
Summary: test
Description: test
Author: Schoperation
"""

load("render.star", "render")
load("schema.star", "schema")
load("http.star", "http")
load("encoding/json.star", "json")
load("encoding/csv.star", "csv")
load("encoding/base64.star", "base64")

def main(config):
    lang_cfg = config.str("lang", "en-US")
    anime_cfg = config.str("anime", "spy-x-family") # Heh
    title_color_cfg = config.str("title_color", "#ffc266")
    sub_id_color_cfg = config.str("sub_id_color", "#a6a6a6")
    dub_id_color_cfg = config.str("dub_id_color", "#a6a6a6")

    anime_atlas_url = "https://raw.githubusercontent.com/Schoperation/Crunchyroll-Anime-Checker/master/lists/anime_atlas_{}.json".format(lang_cfg)
    resp = http.get(url = anime_atlas_url, headers = {"Accept": "application/json", "User-Agent": "Tidbyt - Crunchyroll Anime Status - Schoperation"}, ttl_seconds = 300)
    if resp.status_code != 200:
        print("Failed to load anime atlas with status %d", resp.status_code)
        return render.Root(
            child = render.WrappedText(
                content = "Couldn't get anime :("
            )
        )

    anime_atlas = json.decode(resp.body())["anime"]
    anime = anime_atlas[anime_cfg]

    image_cfg = config.str("anime_image", "poster_full")
    anime_image = ""
    if image_cfg == "sub_thumb":
        anime_image = anime["sub"]["thumbnail"]
    elif image_cfg == "dub_thumb":
        anime_image = anime["dub"]["thumbnail"]
    else:
        anime_image = anime["poster"]

    anime_image = base64.decode(anime_image)

    return render.Root(
        child = render.Column(
            children = [
                render.Marquee(
                    width = 64,
                    align = "center",
                    child = render.Text(
                        font = "CG-pixel-3x5-mono",
                        color = title_color_cfg,
                        content = anime["name"],
                    ),
                ),
                render.Row(
                    children = [
                        render.Image( 
                            width = 27, 
                            height = 27,
                            src = anime_image,
                        ),
                        render.Column(
                            children = [
                                render.Marquee(
                                    width = 36,
                                    align = "center",
                                    child = render.Text(
                                        font = "CG-pixel-3x5-mono",
                                        content = "",
                                    ),
                                ),
                                render.Marquee(
                                    width = 36,
                                    align = "center",
                                    child = render.Text(
                                        font = "CG-pixel-3x5-mono",
                                        color = sub_id_color_cfg,
                                        content = "S:S{s}E{e}".format(s=anime["sub"]["season"], e=anime["sub"]["episode"]),
                                    ),
                                ),
                                render.Marquee(
                                    width = 36,
                                    align = "start",
                                    offset_start = 10,
                                    child = render.Text(
                                        font = "CG-pixel-3x5-mono",
                                        content = anime["sub"]["title"],
                                    ),
                                ),
                                render.Marquee(
                                    width = 36,
                                    align = "center",
                                    child = render.Text(
                                        font = "CG-pixel-3x5-mono",
                                        color = dub_id_color_cfg,
                                        content = "D:S{s}E{e}".format(s=anime["dub"]["season"], e=anime["dub"]["episode"]),
                                    ),
                                ),
                                render.Marquee(
                                    width = 36,
                                    align = "start",
                                    offset_start = 10,
                                    child = render.Text(
                                        font = "CG-pixel-3x5-mono",
                                        content = anime["dub"]["title"],
                                    ),
                                ),
                            ],
                        ),
                    ],
                ),
            ]
        ),
    )

# Temporary while I don't have a csv up
# TODO remove later
CSV_STRING = """
series_id,slug,name
G4PH0WXVJ,spy-x-family,Spy x Family
"""

def get_schema():
    anime_csv = csv.read_all(source = CSV_STRING, skip = 1)
    anime_options = [anime_to_schema_option(anime) for anime in anime_csv]

    config_fields = [
        schema.Dropdown(
            id = "lang",
            name = "Language",
            desc = "Language of subs and dubs to search for.",
            icon = "language",
            default = "en-US",
            options = [
                schema.Option(
                    display = "English (US)",
                    value = "en-US"
                ),
            ],
        ),
        schema.Dropdown(
            id = "anime",
            name = "Anime",
            desc = "The anime you want to check!",
            icon = "tv",
            default = "spy-x-family", # Heh
            options = anime_options,
        ),
        schema.Dropdown(
            id = "anime_image",
            name = "Image",
            desc = "The image to show besides the info.",
            icon = "image",
            default = "poster_full",
            options = [
                schema.Option(
                    display = "Poster (Full)",
                    value = "poster_full",
                ),
                schema.Option(
                    display = "Latest Episode Thumbnail (Sub)",
                    value = "sub_thumb",
                ),
                schema.Option(
                    display = "Latest Episode Thumbnail (Dub)",
                    value = "dub_thumb"
                )
            ]
        ),
        schema.Color(
            id = "title_color",
            name = "Anime Title Color",
            desc = "Color of the anime's title at the top.",
            icon = "brush",
            default = "#ffc266"
        ),
        schema.Color(
            id = "sub_id_color",
            name = "Sub Identifier Color",
            desc = "Color of the latest sub's identifier (S:S1E2).",
            icon = "brush",
            default = "#a6a6a6"
        ),
        schema.Color(
            id = "dub_id_color",
            name = "Dub Identifier Color",
            desc = "Color of the latest dub's identifier (D:S1E2).",
            icon = "brush",
            default = "#a6a6a6"
        ),
    ]

    return schema.Schema(
        version = "1",
        fields = config_fields,
    )

def anime_to_schema_option(anime):
    return schema.Option(
        display = anime[2],
        value = anime[1],
    )
