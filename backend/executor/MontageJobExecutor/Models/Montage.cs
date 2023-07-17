using Newtonsoft.Json;

namespace MontageJobExecutor.Models
{
    public class Montage
    {
        [JsonProperty("ID")]
        public string Id { get; set; }

        [JsonProperty("JobID")]
        public string JobId { get; set; }

        public string Name { get; set; }

        public string FilePath { get; set; }
    }
}
