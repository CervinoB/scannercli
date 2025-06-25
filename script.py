import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import glob
import os
import re
from datetime import datetime

# Set style for better looking plots
# sns.set_style("whitegrid")
plt.rcParams["figure.figsize"] = (16, 8)

# # Configuration
path = "data/iola"
output_dir = os.path.join(path, "visualizations")
os.makedirs(output_dir, exist_ok=True)

# # Load all CSVs
# all_files = glob.glob(os.path.join(path, "issues-*.csv"))
all_files = glob.glob(os.path.join(path, "**", "issues-*.csv"), recursive=True)
df = pd.concat((pd.read_csv(f) for f in all_files), ignore_index=True)


def extract_package(component):
    try:
        return component.split(":")[1].split("/")[1]
    except (IndexError, AttributeError):
        return "unknown"


df["Package"] = df["Component"].apply(extract_package)


# Convert effort to numeric (handling 'h' and 'min')
def convert_effort(effort_str):
    if pd.isna(effort_str):
        return 0
    effort_str = str(effort_str).strip().lower()

    # Extract hours and minutes using regex
    hours = re.search(r"(\d+)\s*h", effort_str)
    minutes = re.search(r"(\d+)\s*min", effort_str)

    total_minutes = 0
    if hours:
        total_minutes += int(hours.group(1)) * 60
    if minutes:
        total_minutes += int(minutes.group(1))

    return total_minutes


df["EffortMinutes"] = df["Effort"].apply(convert_effort)

# Severity counts - safe reindexing
severity_order = ["BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"]
severity_counts = df["Severity"].value_counts()
severity_counts = severity_counts.reindex(severity_order).fillna(0)

# 1. Issues by Severity
# plt.figure(figsize=(10, 6))
# ax = severity_counts.plot(kind="bar", color=sns.color_palette("YlOrRd", n_colors=5))
# plt.title("Issues by Severity Level", pad=20)
# plt.xlabel("Severity Level")
# plt.ylabel("Number of Issues")
# for p in ax.patches:
#     ax.annotate(
#         f"{int(p.get_height())}",
#         (p.get_x() + p.get_width() / 2.0, p.get_height()),
#         ha="center",
#         va="bottom",
#         fontsize=9,
#         xytext=(0, 5),
#         textcoords="offset points",
#     )
# plt.tight_layout()
# plt.savefig(os.path.join(output_dir, "1_issues_by_severity.png"))
# plt.close()


# 3. Issues by Software Quality Dimension

# 4. Top 10 Files with Most Issues
file_issues = df["Component"].value_counts().head(10)
# print(file_issues)
ax = file_issues.plot(kind="barh", color=sns.color_palette("viridis"))
plt.title("Top 10 Files with Most Issues", pad=20)
for p in ax.patches:
    ax.annotate(
        f"{int(p.get_width())}",
        (p.get_width(), p.get_y() + p.get_height() / 2.0),
        ha="left",
        va="center",
        xytext=(5, 0),
        textcoords="offset points",
    )
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "4_top_files_with_issues.png"))
plt.close()

# 5. Effort to Fix by Severity
# effort_by_severity = (
#     df.groupby("Severity")["EffortMinutes"].sum().reindex(severity_order).fillna(0)
# )
# ax = effort_by_severity.plot(kind="bar", color=sns.color_palette("YlOrRd", n_colors=5))
# plt.title("Total Fixing Effort by Severity (Minutes)", pad=20)
# plt.xlabel("Severity Level")
# for p in ax.patches:
#     ax.annotate(
#         f"{int(p.get_height())}",
#         (p.get_x() + p.get_width() / 2.0, p.get_height()),
#         ha="center",
#         va="bottom",
#         xytext=(0, 5),
#         textcoords="offset points",
#     )
# # plt.tight_layout()
# plt.savefig(os.path.join(output_dir, "5_effort_by_severity.png"))
# plt.close()

# 6. Top Issue Messages
top_messages = df["Message"].value_counts().head(20)
print(top_messages)
ax = top_messages.plot(kind="barh", color=sns.color_palette("cubehelix"))
plt.title("Top 15 Most Common Issue Messages", pad=20)
for p in ax.patches:
    ax.annotate(
        f"{int(p.get_width())}",
        (p.get_width(), p.get_y() + p.get_height() / 2.0),
        ha="left",
        va="center",
        xytext=(5, 0),
        textcoords="offset points",
    )
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "6_top_issue_messages.png"))
plt.close()

# 7. Top Packages with Most Issues
package_issues = df["Package"].value_counts().head(10)
ax = package_issues.plot(kind="bar", color=sns.color_palette("husl", n_colors=10))
plt.title("Top 10 Packages with Most Issues", pad=20)
plt.xlabel("Package")
plt.ylabel("Number of Issues")
plt.xticks(rotation=45, ha="right")
for p in ax.patches:
    ax.annotate(
        f"{int(p.get_height())}",
        (p.get_x() + p.get_width() / 2.0, p.get_height()),
        ha="center",
        va="bottom",
        xytext=(0, 5),
        textcoords="offset points",
    )
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "7_issues_by_package.png"))
plt.close()


# Recursively collect all CSVs from subproject folders


# Build combined dataframe with filename metadata
records = []
for file in all_files:
    df = pd.read_csv(file)
    df["SourceFile"] = os.path.basename(file)
    df["Project"] = os.path.basename(os.path.dirname(file))

    # Attempt to parse version/date from filename
    version_part = os.path.basename(file).replace("issues-", "").replace(".csv", "")
    df["Version"] = version_part
    records.append(df)

df_all = pd.concat(records, ignore_index=True)

# messages = df_all["Message"].value_counts().head(20)
messages = [
    "Extract this nested ternary operation into an independent statement.",
    "Refactor this code to not nest functions more than 4 levels deep.",
    "Handle this exception or don't catch it at all.",
]
messages = [m.upper() for m in messages]
print(f"Unique messages found: {messages}")

# Standardize severity categories
# severity_order = ["BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"]
df_all["Message"] = df_all["Message"].fillna("UNKNOWN").str.upper()
df_all["Message"] = df_all["Message"].where(
    df_all["Message"].isin(messages),
    "Refactor this code to not nest functions more than 4 levels deep.",
)

# Aggregate severity counts by Version
agg_df = df_all.groupby(["Version", "Message"]).size().unstack(fill_value=0)
agg_df = agg_df[messages]  # reorder columns

agg_df.index = agg_df.index.astype(str)


# Then apply the natural sort
def natural_sort_key(version_str):
    # Split into numeric and non-numeric parts, preserving order
    parts = re.split(r"(\d+)", version_str.replace("v", ""))
    return [int(p) if p.isdigit() else p.lower() for p in parts]


agg_df = agg_df.loc[sorted(agg_df.index, key=natural_sort_key)]
print(agg_df)


# Plot time series stacked area chart
# plt.figure(figsize=(14, 8))
ax = agg_df.plot.area(stacked=True, cmap="Spectral", alpha=0.9, linewidth=0)
plt.title("Série Temporal de Alertas do SonarQube por PRs", fontsize=14, pad=20)
plt.xlabel("Versão ou Data do Scan")
plt.ylabel("Número de Alertas")
plt.legend(title="Mensagem", loc="upper left")
plt.xticks(rotation=45, ha="right")
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "8_time_series_message_area.png"))
plt.close()


# Count occurrences of the specific message per version
target_message = "Refactor this code to not nest functions more than 4 levels deep."
message_counts = df_all[df_all["Message"] == target_message].groupby("Version").size()
message_counts = message_counts.loc[sorted(message_counts.index, key=natural_sort_key)]


# Plot bar chart of the message count per version
ax = message_counts.plot(kind="bar", color="salmon")
plt.title("Occurrences of Nesting Warning by Version", pad=20)
plt.xlabel("Version")
plt.ylabel("Number of Occurrences")
plt.xticks(rotation=45, ha="right")
for p in ax.patches:
    ax.annotate(
        f"{int(p.get_height())}",
        (p.get_x() + p.get_width() / 2.0, p.get_height()),
        ha="left",
        va="bottom",
        xytext=(0, 5),
        textcoords="offset points",
    )
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "9_nesting_warning_by_version.png"))
plt.close()

print(f"Visualizations saved to: {output_dir}")
