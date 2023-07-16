using MontageJobExecutor.Contracts;
using OpenAI.ObjectModels;
using OpenAI.Interfaces;

namespace MontageJobExecutor.Services
{
    public class WhisperVideoTranscriber : IVideoTranscriber
    {
        readonly IOpenAIService _openAIService;

        public WhisperVideoTranscriber(IOpenAIService openAIService)
        {
            _openAIService = openAIService;
        }

        public async Task<TranscribeResponse> TranscribeAsync(TranscribeRequest request)
        {
            var fileName = Path.GetFileName(request.SourcePath);
            var outputPath = Path.Combine(
                Path.GetDirectoryName(request.SourcePath) ?? "",
                $"{Path.GetFileNameWithoutExtension(request.SourcePath)}.srt"
            );

            var result = await _openAIService.Audio.CreateTranscription(new()
            {
                FileName       = fileName,
                FileStream     = File.OpenRead(request.SourcePath),
                Model          = OpenAI.ObjectModels.Models.WhisperV1,
                ResponseFormat = StaticValues.AudioStatics.ResponseFormat.Srt,
            });
            if (!result.Successful || result.Error != null)
                throw new Exception($"Failed to transcribe video at path {request.SourcePath} . Error: {result.Error}");

            Console.WriteLine($"Transcription for video at path {request.SourcePath} took {result.Duration} . Output at {outputPath}");
            File.WriteAllText(outputPath, result.Text);
            
            return new()
            {
                OutputPath = outputPath,
            };
        }
    }
}
