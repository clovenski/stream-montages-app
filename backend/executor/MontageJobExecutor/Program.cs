using Confluent.Kafka;
using MontageJobExecutor.Models;
using MontageJobExecutor.Services;
using OpenAI;
using OpenAI.Managers;
using System.Text.Json;

const string ENV_KAFKA_CONSUMER_BOOTSTRAP_SERVERS = "SM_KAFKA_CONSUMER_BOOTSTRAP_SERVERS";
const string ENV_KAFKA_CONSUMER_GROUP_ID          = "SM_KAFKA_CONSUMER_GROUP_ID";
const string ENV_KAFKA_TOPIC                      = "SM_KAFKA_TOPIC";
const string ENV_JOB_REPO_BASE_URL                = "SM_JOB_REPO_BASE_URL";
const string ENV_YOUTUBE_DL_EXE_PATH              = "SM_YOUTUBE_DL_EXE_PATH";
const string ENV_OPENAI_API_KEY                   = "SM_OPENAI_API_KEY";
const string ENV_MONTAGE_OUTPUT_PATH_BASE         = "SM_MONTAGE_OUTPUT_PATH_BASE";

var config = new ConsumerConfig
{
    BootstrapServers = Environment.GetEnvironmentVariable(ENV_KAFKA_CONSUMER_BOOTSTRAP_SERVERS),
    GroupId          = Environment.GetEnvironmentVariable(ENV_KAFKA_CONSUMER_GROUP_ID),
    AutoOffsetReset  = AutoOffsetReset.Earliest,
};

using var consumer = new ConsumerBuilder<Ignore, string>(config).Build();

consumer.Subscribe(Environment.GetEnvironmentVariable(ENV_KAFKA_TOPIC));

var httpClient = new HttpClient();

var jobRepo                 = new MontageJobRepository(httpClient, Environment.GetEnvironmentVariable(ENV_JOB_REPO_BASE_URL)!);
var youtubeAVStreamProvider = new YouTubeAVStreamProvider(Environment.GetEnvironmentVariable(ENV_YOUTUBE_DL_EXE_PATH));
var clipService             = new FFmpegClipService();
var videoTranscriber        = new WhisperVideoTranscriber(new OpenAIService(new OpenAiOptions() { ApiKey = Environment.GetEnvironmentVariable(ENV_OPENAI_API_KEY)! }, httpClient));

var montageBuilder = new MontageBuilder(
   youtubeAVStreamProvider,
   clipService,
   videoTranscriber
);

try
{
    while (true)
    {
        var result = consumer.Consume();
        if (result.IsPartitionEOF || result.Message == null)
            continue;

        var message = result.Message.Value;

        Console.WriteLine($"Got this message: {message}");

        MontageJob? job = null;
        try
        {
           // json parse to job dto
           job = JsonSerializer.Deserialize<MontageJob>(message)
               ?? throw new JsonException("Deserialization resulted in null");

           // update the job to started status
           job.Status = "STARTED";
           job = (await jobRepo.UpdateMontageJobAsync(new() { Entity = job })).UpdatedEntity;
           Console.WriteLine($"Successfully marked job {job.ID} as started.");

           var montageOutputPath = Path.Combine(
               Environment.GetEnvironmentVariable(ENV_MONTAGE_OUTPUT_PATH_BASE)!,
               $"{job.ID}"
           );
           Directory.CreateDirectory(montageOutputPath);
           
           var montage = await montageBuilder.Build(new()
           {
               Highlights     = job.JobDefinition.Highlights,
               OutputPathBase = montageOutputPath,
           });
           Console.WriteLine($"Montage built, video at {montage.VideoPath} transcription at {montage.TranscriptionPath}");

           // mark job as complete
           job.Status = "COMPLETE";
           await jobRepo.UpdateMontageJobAsync(new() { Entity = job });
        }
        catch (Exception ex)
        {
           Console.WriteLine($"Failed to process montage job. {ex}");

           if (job != null)
           {
               job.Status = "FAILED";
               await jobRepo.UpdateMontageJobAsync(new() { Entity = job });
               Console.WriteLine($"Successfully marked job {job.ID} as failed.");
           }
        }
    }
}
catch (Exception ex)
{
    Console.WriteLine($"Unhandled error occurred: {ex}");
}
finally
{
    consumer.Close();
}
