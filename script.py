import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import glob
import os
import re
from datetime import datetime

# Set style for better looking plots
sns.set_style("whitegrid")
plt.rcParams["figure.figsize"] = (12, 6)

# Configuration
path = "data/twenty-aux-bbk/"
output_dir = os.path.join(path, "visualizations")
os.makedirs(output_dir, exist_ok=True)

# Load all CSVs
all_files = glob.glob(os.path.join(path, "issues-*.csv"))
df = pd.concat((pd.read_csv(f) for f in all_files), ignore_index=True)


# Clean and prepare data
# Extract package name from component
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
plt.figure(figsize=(10, 6))
ax = severity_counts.plot(kind="bar", color=sns.color_palette("YlOrRd", n_colors=5))
plt.title("Issues by Severity Level", pad=20)
plt.xlabel("Severity Level")
plt.ylabel("Number of Issues")
for p in ax.patches:
    ax.annotate(
        f"{int(p.get_height())}",
        (p.get_x() + p.get_width() / 2.0, p.get_height()),
        ha="center",
        va="bottom",
        fontsize=9,
        xytext=(0, 5),
        textcoords="offset points",
    )
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "1_issues_by_severity.png"))
plt.close()

# 2. Issues by Type (limit pie chart if many types)
type_counts = df["Type"].value_counts()
if len(type_counts) > 15:
    type_counts = type_counts.head(14)
    type_counts["Other"] = df["Type"].value_counts()[14:].sum()

plt.figure(figsize=(10, 6))
ax = type_counts.plot(
    kind="pie",
    autopct="%1.1f%%",
    startangle=90,
    colors=sns.color_palette("pastel", n_colors=len(type_counts)),
    wedgeprops={"linewidth": 1, "edgecolor": "white"},
)
plt.title("Distribution of Issue Types", pad=20)
plt.ylabel("")
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "2_issues_by_type.png"))
plt.close()

# 3. Issues by Software Quality Dimension
quality_counts = df["SoftwareQuality"].value_counts()
plt.figure(figsize=(12, 6))
ax = quality_counts.plot(kind="barh", color=sns.color_palette("Blues_d"))
plt.title("Issues by Software Quality Dimension", pad=20)
plt.xlabel("Number of Issues")
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
plt.savefig(os.path.join(output_dir, "3_issues_by_quality_dimension.png"))
plt.close()

# 4. Top 10 Files with Most Issues
plt.figure(figsize=(12, 8))
file_issues = df["Component"].value_counts().head(10)
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
effort_by_severity = (
    df.groupby("Severity")["EffortMinutes"].sum().reindex(severity_order).fillna(0)
)
plt.figure(figsize=(10, 6))
ax = effort_by_severity.plot(kind="bar", color=sns.color_palette("YlOrRd", n_colors=5))
plt.title("Total Fixing Effort by Severity (Minutes)", pad=20)
plt.xlabel("Severity Level")
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
plt.savefig(os.path.join(output_dir, "5_effort_by_severity.png"))
plt.close()

# 6. Top Issue Messages
top_messages = df["Message"].value_counts().head(15)
plt.figure(figsize=(12, 8))
ax = top_messages.plot(kind="barh", color=sns.color_palette("cubehelix", n_colors=15))
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
plt.figure(figsize=(12, 8))
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
all_files = glob.glob(os.path.join(path, "**", "issues-*.csv"), recursive=True)

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

# Standardize severity categories
severity_order = ["BLOCKER", "CRITICAL", "MAJOR", "MINOR", "INFO"]
df_all["Severity"] = df_all["Severity"].fillna("UNKNOWN").str.upper()
df_all["Severity"] = df_all["Severity"].where(
    df_all["Severity"].isin(severity_order), "INFO"
)

# Aggregate severity counts by Version
agg_df = df_all.groupby(["Version", "Severity"]).size().unstack(fill_value=0)
agg_df = agg_df[severity_order]  # reorder columns

# Sort versions naturally if possible
agg_df = agg_df.sort_index(
    key=lambda x: x.map(
        lambda s: [
            int(n) if n.isdigit() else n
            for n in s.replace("v", "").replace(".", " ").split()
        ]
    )
)

# Plot time series stacked area chart
plt.figure(figsize=(14, 8))
ax = agg_df.plot.area(stacked=True, cmap="Spectral", alpha=0.9, linewidth=0)
plt.title("Série Temporal de Alertas do SonarQube por PRs", fontsize=14, pad=20)
plt.xlabel("Versão ou Data do Scan")
plt.ylabel("Número de Alertas")
plt.legend(title="Severidade", loc="upper left")
plt.xticks(rotation=45, ha="right")
plt.tight_layout()
plt.savefig(os.path.join(output_dir, "8_time_series_severity_area.png"))
plt.close()

print(f"Visualizations saved to: {output_dir}")
