using System.Text.Json.Serialization;

namespace MontageJobExecutor.Models
{
    public class MontageJob
    {
        public string ID { get; set; }

        public string Name { get; set; }

        public string Status { get; set; }

        public string CreatedAt { get; set; }

        public string UpdatedAt { get; set; }

        public MontageJobDefinition JobDefinition { get; set; }
    }

    public class MontageJobDefinition
    {
        public List<HighlightInfo> Highlights { get; set; }
    }

    public class HighlightInfo
    {
        public string VideoURL { get; set; }

        /// <summary>
        /// Ex. 1:20 means the highlight starts at 1 hour, 20 minutes
        /// </summary>
        [JsonPropertyName("Timestamp")]
        public string MinuteTimestamp { get; set; }

        public int DurationSeconds { get; set; }
    }
}
