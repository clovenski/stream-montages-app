using MontageJobExecutor.Contracts;
using MontageJobExecutor.Models;
using System.Net.Http.Json;
using System.Text.Json;

namespace MontageJobExecutor.Services
{
    public class MontageRepository : IMontageRepository
    {
        readonly HttpClient _httpClient;
        readonly string _serviceBaseUrl;

        public MontageRepository(HttpClient httpClient, string serviceBaseUrl)
        {
            _httpClient     = httpClient;
            _serviceBaseUrl = serviceBaseUrl.TrimEnd('/');
        }

        public async Task<CreateMontageResponse> CreateMontageAsync(CreateMontageRequest request)
        {
            var response = await _httpClient.PostAsJsonAsync($"{_serviceBaseUrl}/montages", request.Entity);

            response.EnsureSuccessStatusCode();

            return new()
            {
                CreatedEntity = await JsonSerializer.DeserializeAsync<Montage>(response.Content.ReadAsStream()),
            };
        }
    }
}
