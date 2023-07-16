using MontageJobExecutor.Contracts;
using MontageJobExecutor.Models;
using System.Net.Http.Json;
using System.Text.Json;

namespace MontageJobExecutor.Services
{
    public class MontageJobRepository : IMontageJobRepository
    {
        readonly HttpClient _httpClient;
        readonly string _serviceBaseUrl;

        public MontageJobRepository(HttpClient httpClient, string serviceBaseUrl)
        {
            _httpClient     = httpClient;
            _serviceBaseUrl = serviceBaseUrl.TrimEnd('/');
        }

        public async Task<UpdateMontageJobResponse> UpdateMontageJobAsync(UpdateMontageJobRequest request)
        {
            var response = await _httpClient.PutAsJsonAsync($"{_serviceBaseUrl}/montages/jobs/{request.Entity.ID}", request.Entity);

            response.EnsureSuccessStatusCode();

            return new()
            {
                UpdatedEntity = await JsonSerializer.DeserializeAsync<MontageJob>(response.Content.ReadAsStream()),
            };
        }
    }
}
